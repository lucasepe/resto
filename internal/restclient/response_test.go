package restclient

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDumpResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantOK     string
		wantKO     string
		wantErr    bool
	}{
		{
			name:       "success 200",
			statusCode: 200,
			body:       "hello success",
			wantOK:     "hello success",
			wantKO:     "",
			wantErr:    false,
		},
		{
			name:       "error 404",
			statusCode: 404,
			body:       "not found",
			wantOK:     "",
			wantKO:     "not found",
			wantErr:    true,
		},
		{
			name:       "no body success",
			statusCode: 204,
			body:       "",
			wantOK:     "",
			wantKO:     "",
			wantErr:    false,
		},
		{
			name:       "no body error",
			statusCode: 500,
			body:       "",
			wantOK:     "",
			wantKO:     "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			var okBuf, koBuf bytes.Buffer
			err := dumpResponse(resp, &okBuf, &koBuf)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error status: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if got := okBuf.String(); got != tt.wantOK {
				t.Errorf("okWri got %q, want %q", got, tt.wantOK)
			}
			if got := koBuf.String(); got != tt.wantKO {
				t.Errorf("koWri got %q, want %q", got, tt.wantKO)
			}
		})
	}
}

func TestIsTextResponse(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expectText  bool
	}{
		{"empty content type", "", true},
		{"text/plain", "text/plain", true},
		{"text/html", "text/html", true},
		{"application/json", "application/json", false},
		{"invalid content type", "invalid/type", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{"Content-Type": []string{tc.contentType}}}
			if isTextResponse(resp) != tc.expectText {
				t.Errorf("unexpected result for %s: got %v, want %v", tc.contentType, !tc.expectText, tc.expectText)
			}
		})
	}
}
