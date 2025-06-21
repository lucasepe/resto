package restclient

import "testing"

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		hasCA    bool
		hasBasic bool
		hasToken bool
		hasCert  bool
	}{
		{
			name:   "Empty Config",
			config: Config{},
			hasCA:  false, hasBasic: false, hasToken: false, hasCert: false,
		},
		{
			name:   "Has CA Data",
			config: Config{CertificateAuthorityData: "cert-data"},
			hasCA:  true, hasBasic: false, hasToken: false, hasCert: false,
		},
		{
			name:   "Has Basic Auth",
			config: Config{Username: "user", Password: "pass"},
			hasCA:  false, hasBasic: true, hasToken: false, hasCert: false,
		},
		{
			name:   "Has Token Auth",
			config: Config{Token: "my-token"},
			hasCA:  false, hasBasic: false, hasToken: true, hasCert: false,
		},
		{
			name:   "Has Certificate Auth",
			config: Config{ClientCertificateData: "cert", ClientKeyData: "key"},
			hasCA:  false, hasBasic: false, hasToken: false, hasCert: true,
		},
		{
			name: "Has All Auth Methods",
			config: Config{
				CertificateAuthorityData: "cert-data",
				Username:                 "user",
				Password:                 "pass",
				Token:                    "my-token",
				ClientCertificateData:    "cert",
				ClientKeyData:            "key",
			},
			hasCA: true, hasBasic: true, hasToken: true, hasCert: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.HasCA(); got != tt.hasCA {
				t.Errorf("HasCA() = %v, want %v", got, tt.hasCA)
			}
			if got := tt.config.HasBasicAuth(); got != tt.hasBasic {
				t.Errorf("HasBasicAuth() = %v, want %v", got, tt.hasBasic)
			}
			if got := tt.config.HasTokenAuth(); got != tt.hasToken {
				t.Errorf("HasTokenAuth() = %v, want %v", got, tt.hasToken)
			}
			if got := tt.config.HasCertAuth(); got != tt.hasCert {
				t.Errorf("HasCertAuth() = %v, want %v", got, tt.hasCert)
			}
		})
	}
}
