package restclient

import (
	"fmt"
	"log"
	"net/http"
)

func HTTPClientForConfig(cfg Config) (*http.Client, error) {
	rt, err := tlsConfigFor(&cfg)
	if err != nil {
		return &http.Client{
			Transport: defaultTransport(),
		}, err
	}

	if cfg.Verbose {
		log.Println("using verbose roundtripper")

		rt = &verboseRoundTripper{
			v:    true,
			next: rt,
		}
	}

	// Set authentication wrappers
	switch {
	case cfg.HasBasicAuth() && cfg.HasTokenAuth():
		return nil, fmt.Errorf("username/password or bearer token may be set, but not both")

	case cfg.HasTokenAuth():
		if cfg.Verbose {
			log.Println("using bearer auth roundtripper")
		}
		rt = &bearerAuthRoundTripper{
			bearer: cfg.Token,
			next:   rt,
		}

	case cfg.HasBasicAuth():
		if cfg.Verbose {
			log.Println("using basic auth roundtripper")
		}
		rt = &basicAuthRoundTripper{
			username: cfg.Username,
			password: cfg.Password,
			next:     rt,
		}
	}

	return &http.Client{Transport: rt}, nil
}
