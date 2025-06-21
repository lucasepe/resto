package restclient

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProxyURL(t *testing.T) {
	tests := []struct {
		name      string
		proxyURL  string
		expectErr bool
	}{
		{"Valid HTTP proxy", "http://proxy.example.com", false},
		{"Valid HTTPS proxy", "https://secure-proxy.example.com", false},
		{"Invalid scheme", "ftp://invalid-proxy.com", true},
		{"Malformed URL", ":://bad-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := parseProxyURL(tt.proxyURL)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, u)
			}
		})
	}
}

func TestBasicAuthRoundTripper(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte("user:pass")), r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	})

	mockRT := &mockRoundTripper{mux: mux}

	client := &basicAuthRoundTripper{
		username: "user",
		password: "pass",
		next:     mockRT,
	}

	req, _ := http.NewRequest("GET", "https://httpbin.org", nil)
	_, err := client.RoundTrip(req)
	assert.NoError(t, err)
}

func TestBearerAuthRoundTripper(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer mytoken", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	})

	mockRT := &mockRoundTripper{mux: mux}

	client := &bearerAuthRoundTripper{
		bearer: "mytoken",
		next:   mockRT,
	}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := client.RoundTrip(req)
	assert.NoError(t, err)
}

type mockRoundTripper struct {
	mux *http.ServeMux
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	m.mux.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func TestHTTPSCallWithInsecureSkipVerify(t *testing.T) {
	_, _, tlsCert := generateSelfSignedCert()

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
	}))

	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
	server.StartTLS()
	defer server.Close()

	cfg := Config{
		Verbose:   false,
		Insecure:  true,
		ServerURL: server.URL,
	}

	client, err := HTTPClientForConfig(cfg)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, cfg.ServerURL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err, "HTTPS request with Insecure=true should not fail with TLS error")
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(data), "success")
}

func TestHTTPSCallFailsWithoutInsecureSkipVerify(t *testing.T) {
	_, _, tlsCert := generateSelfSignedCert()

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message": "success"}`)
	}))

	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
	server.StartTLS()
	defer server.Close()

	cfg := Config{
		Verbose:   false,
		Insecure:  false,
		ServerURL: server.URL,
	}

	client, err := HTTPClientForConfig(cfg)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, cfg.ServerURL, nil)
	require.NoError(t, err)

	_, err = client.Do(req)
	require.Error(t, err, "Expected TLS error due to unknown certificate authority")

	require.Contains(t, err.Error(), "x509: certificate signed by unknown authority")
}

func TestHTTPSCallSucceedsWithoutInsecureSkipVerify(t *testing.T) {
	certPEM, keyPEM, tlsCert := generateSelfSignedCert()

	caDataBase64 := base64.StdEncoding.EncodeToString(certPEM)
	clientCertBase64 := base64.StdEncoding.EncodeToString(certPEM)
	clientKeyBase64 := base64.StdEncoding.EncodeToString(keyPEM)

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message":"success"}`))
	}))
	server.TLS = &tls.Config{Certificates: []tls.Certificate{tlsCert}}
	server.StartTLS()
	defer server.Close()

	cfg := Config{
		Verbose:                  false,
		Insecure:                 false,
		ServerURL:                server.URL,
		CertificateAuthorityData: caDataBase64,
		ClientCertificateData:    clientCertBase64,
		ClientKeyData:            clientKeyBase64,
	}

	client, err := HTTPClientForConfig(cfg)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, cfg.ServerURL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func generateSelfSignedCert() (certPEM []byte, keyPEM []byte, tlsCert tls.Certificate) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("unable to generate private key: %v", err)
	}

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	// Self-signed cert with SAN: localhost
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		DNSNames: []string{"localhost"},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
	}

	// Signature
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("unable to create certificate: %v", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// Server tls.Certificate
	tlsCert, err = tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("unable to load X509 key pair: %v", err)
	}

	return certPEM, keyPEM, tlsCert
}
