package call

import (
	"reflect"
	"testing"
)

func TestReverseURL(t *testing.T) {
	tests := []struct {
		name       string
		rawurl     string
		wantBase   string
		wantPath   string
		wantParams []string
		wantErr    bool
	}{
		{
			name:       "basic url with params",
			rawurl:     "https://example.com/path?a=1&b=2",
			wantBase:   "https://example.com",
			wantPath:   "/path",
			wantParams: []string{"a:1", "b:2"},
			wantErr:    false,
		},
		{
			name:       "url with multiple values for same param",
			rawurl:     "http://test.org/foo?x=1&x=2&y=3",
			wantBase:   "http://test.org",
			wantPath:   "/foo",
			wantParams: []string{"x:1", "x:2", "y:3"},
			wantErr:    false,
		},
		{
			name:       "url with empty param value",
			rawurl:     "https://host/path?empty=",
			wantBase:   "https://host",
			wantPath:   "/path",
			wantParams: []string{"empty:"},
			wantErr:    false,
		},
		{
			name:       "url without path",
			rawurl:     "https://domain.com",
			wantBase:   "https://domain.com",
			wantPath:   "/",
			wantParams: nil,
			wantErr:    false,
		},
		{
			name:     "url missing scheme",
			rawurl:   "domain.com/path",
			wantPath: "domain.com/path",
			wantErr:  false,
		},
		{
			name:    "invalid url",
			rawurl:  "http://%gh&%ij",
			wantErr: true,
		},
		{
			name:       "url with param without value",
			rawurl:     "https://site.com/path?foo",
			wantBase:   "https://site.com",
			wantPath:   "/path",
			wantParams: []string{"foo:"},
			wantErr:    false,
		},
		{
			name:       "url with empty param key ignored",
			rawurl:     "https://site.com/path?=novalue",
			wantBase:   "https://site.com",
			wantPath:   "/path",
			wantParams: []string{":novalue"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBase, gotPath, gotParams, err := reverseURL(tt.rawurl)
			if (err != nil) != tt.wantErr {
				t.Fatalf("reverseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if gotBase != tt.wantBase {
				t.Errorf("BaseURL = %q, want %q", gotBase, tt.wantBase)
			}
			if gotPath != tt.wantPath {
				t.Errorf("Path = %q, want %q", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("Params = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}
