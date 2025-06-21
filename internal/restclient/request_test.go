package restclient

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestCloneRequest(t *testing.T) {
	orig := &http.Request{
		Header: http.Header{
			"Test-Header": []string{"value1", "value2"},
		},
	}
	clone := cloneRequest(orig)

	if &orig == &clone {
		t.Errorf("CloneRequest did not create a new request instance")
	}
	if &orig.Header == &clone.Header {
		t.Errorf("CloneRequest did not create a deep copy of headers")
	}
	if len(clone.Header.Get("Test-Header")) == 0 {
		t.Errorf("CloneRequest did not copy headers correctly")
	}
}

func TestCloneHeader(t *testing.T) {
	orig := http.Header{
		"Test-Header": []string{"value1", "value2"},
	}
	clone := cloneHeader(orig)

	if orig.Get("Test-Header") != clone.Get("Test-Header") {
		t.Errorf("cloneHeader did not create a deep copy of header values")
	}
}

func TestComposeURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		params   []string
		expected string
		wantErr  bool
	}{
		{
			name:     "base only",
			baseURL:  "https://example.com",
			path:     "",
			params:   nil,
			expected: "https://example.com",
		},
		{
			name:     "base and path",
			baseURL:  "https://example.com/",
			path:     "/api/v1",
			params:   nil,
			expected: "https://example.com/api/v1",
		},
		{
			name:     "with single param",
			baseURL:  "https://example.com",
			path:     "api",
			params:   []string{"name:john"},
			expected: "https://example.com/api?name=john",
		},
		{
			name:     "with multiple params",
			baseURL:  "https://example.com",
			path:     "api",
			params:   []string{"name:john", "role:admin"},
			expected: "https://example.com/api?name=john&role=admin",
		},
		{
			name:     "ignore invalid param (missing colon)",
			baseURL:  "https://example.com",
			path:     "api",
			params:   []string{"invalidparam", "valid:yes"},
			expected: "https://example.com/api?valid=yes",
		},
		{
			name:     "ignore param with empty key",
			baseURL:  "https://example.com",
			path:     "api",
			params:   []string{":value", "a:b"},
			expected: "https://example.com/api?a=b",
		},
		{
			name:     "with empty value",
			baseURL:  "https://example.com",
			path:     "api",
			params:   []string{"flag:"},
			expected: "https://example.com/api?flag=",
		},
		{
			name:    "invalid URL",
			baseURL: "://bad",
			path:    "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := composeURL(tt.baseURL, tt.path, tt.params...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("got URL: %q, want: %q", got, tt.expected)
			}
		})
	}
}

func TestSetHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headers     []string
		wantHeaders map[string]string
		wantLogs    []string
	}{
		{
			name:    "Valid headers",
			headers: []string{"Content-Type: application/json", "X-Custom-Header: value"},
			wantHeaders: map[string]string{
				"Content-Type":    "application/json",
				"X-Custom-Header": "value",
			},
		},
		{
			name:    "Header with empty value",
			headers: []string{"X-Empty-Value:"},
			wantHeaders: map[string]string{
				"X-Empty-Value": "",
			},
		},
		{
			name:     "Invalid header no colon",
			headers:  []string{"InvalidHeader"},
			wantLogs: []string{`ignoring invalid header format: "InvalidHeader"`},
		},
		{
			name:     "Empty key",
			headers:  []string{":value"},
			wantLogs: []string{`ignoring header with empty key: ":value"`},
		},
		{
			name:    "Trim spaces around key and value",
			headers: []string{"  X-Trim  :   spaced value  "},
			wantHeaders: map[string]string{
				"X-Trim": "spaced value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)

			// Catturiamo output stderr (per verificare i log di warning)
			var stderr strings.Builder
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			done := make(chan struct{})
			go func() {
				buf, err := io.ReadAll(r)
				if err == nil {
					stderr.Write(buf)
				}
				close(done)
			}()

			setHeaders(req, tt.headers...)

			w.Close()
			os.Stderr = oldStderr
			<-done

			// Verifica headers
			for k, v := range tt.wantHeaders {
				got := req.Header.Get(k)
				if got != v {
					t.Errorf("header %q = %q, want %q", k, got, v)
				}
			}

			// Verifica che non ci siano header non aspettati
			if len(tt.wantHeaders) != len(req.Header) {
				t.Errorf("unexpected number of headers: got %d, want %d", len(req.Header), len(tt.wantHeaders))
			}

			// Verifica log di warning se previsto
			for _, wantLog := range tt.wantLogs {
				if !strings.Contains(stderr.String(), wantLog) {
					t.Errorf("expected stderr to contain %q, got %q", wantLog, stderr.String())
				}
			}
		})
	}
}
