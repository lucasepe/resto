package restclient

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRESTClient_Do(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/testpath" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("method not allowed"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}))
	defer ts.Close()

	cfg := ConfigFromEnv()
	cfg.ServerURL = ts.URL

	cli, err := HTTPClientForConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	outBuf, errBuf := bytes.Buffer{}, bytes.Buffer{}
	streams := IOStreams{
		Out: &outBuf,
		Err: &errBuf,
	}

	err = New(
		RequestOptions{
			BaseURL: cfg.ServerURL,
			Path:    "/testpath",
			Params:  []string{"foo:bar", "baz:qux"},
			Headers: []string{"X-Test-Header: testvalue"},
		},
	).Do(context.Background(), cli, streams)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if got := outBuf.String(); !strings.Contains(got, "hello") {
		t.Errorf("stdout does not contain expected 'hello', got: %q", got)
	}
	if errBuf.Len() != 0 {
		t.Errorf("stderr should be empty, got: %q", errBuf.String())
	}
}
