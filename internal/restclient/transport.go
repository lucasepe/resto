package restclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func tlsConfigFor(ep *Config) (http.RoundTripper, error) {
	res := defaultTransport()

	if ep.ProxyURL != "" {
		u, err := parseProxyURL(ep.ProxyURL)
		if err != nil {
			return nil, err
		}

		res.Proxy = http.ProxyURL(u)
	}

	caCertPool := x509.NewCertPool()

	if len(ep.CertificateAuthorityData) > 0 {
		caData, err := base64.StdEncoding.DecodeString(ep.CertificateAuthorityData)
		if err != nil {
			return nil, fmt.Errorf("unable to decode certificate authority data")
		}

		caCertPool.AppendCertsFromPEM(caData)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: ep.Insecure,
		RootCAs:            caCertPool,
	}
	defer func() {
		res.TLSClientConfig = tlsConfig
	}()

	if !ep.HasCertAuth() {
		return res, nil
	}

	certData, err := base64.StdEncoding.DecodeString(ep.ClientCertificateData)
	if err != nil {
		return nil, fmt.Errorf("unable to decode client certificate data")
	}

	keyData, err := base64.StdEncoding.DecodeString(ep.ClientKeyData)
	if err != nil {
		return nil, fmt.Errorf("unable to decode client key data")
	}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return res, err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	return res, nil
}

func defaultTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func parseProxyURL(proxyURL string) (*url.URL, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %v", proxyURL)
	}

	switch u.Scheme {
	case "http", "https", "socks5":
	default:
		return nil, fmt.Errorf("unsupported scheme %q, must be http, https, or socks5", u.Scheme)
	}
	return u, nil
}

type basicAuthRoundTripper struct {
	username string
	password string `datapolicy:"password"`
	next     http.RoundTripper
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.next.RoundTrip(req)
	}
	req = cloneRequest(req)
	req.SetBasicAuth(rt.username, rt.password)
	return rt.next.RoundTrip(req)
}

type bearerAuthRoundTripper struct {
	bearer string
	next   http.RoundTripper
}

func (rt *bearerAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.next.RoundTrip(req)
	}

	req = cloneRequest(req)
	token := rt.bearer

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return rt.next.RoundTrip(req)
}

type verboseRoundTripper struct {
	v    bool
	next http.RoundTripper
}

func (vt *verboseRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	if vt.v {
		fmt.Fprintln(os.Stderr)

		dumpReq, _ := httputil.DumpRequestOut(req, false)
		addPrefixToLines(os.Stderr, dumpReq, "> ")

		if len(reqBody) > 0 {
			fmt.Fprintf(os.Stderr, "\n%s\n", string(reqBody))
		}
		fmt.Fprint(os.Stderr, "\n\n")
	}

	resp, err := vt.next.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}

	if vt.v {
		dumpResp, _ := httputil.DumpResponse(resp, false)
		addPrefixToLines(os.Stderr, dumpResp, "< ")

		contentType := resp.Header.Get("Content-Type")
		fmt.Fprintln(os.Stderr)
		if strings.Contains(contentType, "application/json") {
			prettyPrintJSON(respBody)
		} else {
			fmt.Fprintln(os.Stderr, string(respBody))
		}
	}

	return resp, nil
}

func prettyPrintJSON(body []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, body, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, string(body))
		return
	}
	fmt.Fprintln(os.Stderr, out.String())
}

func addPrefixToLines(w io.Writer, data []byte, prefix string) {
	lines := bytes.SplitSeq(data, []byte("\n"))
	for line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		fmt.Fprintf(w, "%s%s\n", prefix, line)
	}
}
