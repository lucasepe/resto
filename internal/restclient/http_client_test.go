package restclient

import (
	"testing"
)

func TestHTTPClientForConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		expectErr bool
	}{
		{
			name:      "valid config without auth",
			cfg:       Config{},
			expectErr: false,
		},
		{
			name:      "valid config with bearer token",
			cfg:       Config{Token: "test-token"},
			expectErr: false,
		},
		{
			name:      "valid config with basic auth",
			cfg:       Config{Username: "user", Password: "pass"},
			expectErr: false,
		},
		{
			name:      "invalid config with both auth methods",
			cfg:       Config{Username: "user", Password: "pass", Token: "token"},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := HTTPClientForConfig(tc.cfg)
			if (err != nil) != tc.expectErr {
				t.Errorf("unexpected error status: got %v, expectErr %v", err, tc.expectErr)
			}
			if !tc.expectErr && client == nil {
				t.Errorf("expected client, got nil")
			}
		})
	}
}
