package call

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// reverseURL parses the input URL string and returns:
// - baseURL: scheme + host (e.g. "https://example.com")
// - urlPath: the path component (default to "/" if empty)
// - params: slice of "key:val" strings from query parameters (val may be empty)
// Returns an error if URL parsing fails or if scheme/host are missing.
func reverseURL(rawurl string) (baseURL, urlPath string, params []string, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to parse url: %w", err)
	}

	// If both scheme and host are empty, treat it as a relative path
	if u.Scheme == "" && u.Host == "" {
		baseURL = ""
		urlPath = u.Path
		if urlPath == "" {
			urlPath = "/"
		}
	} else if u.Scheme != "" && u.Host != "" {
		baseURL = u.Scheme + "://" + u.Host
		urlPath = u.Path
		if urlPath == "" {
			urlPath = "/"
		}
	} else {
		// Either scheme or host is missing (incomplete URL)
		return "", "", nil, errors.New("invalid URL: missing scheme or host")
	}

	urlPath = u.Path
	if urlPath == "" {
		urlPath = "/"
	}

	query := u.Query()
	// per avere un output prevedibile ordino le chiavi
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vals := query[k]
		if len(vals) == 0 {
			// param senza valore: key:
			params = append(params, k+":")
			continue
		}
		for _, v := range vals {
			params = append(params, k+":"+v)
		}
	}

	return baseURL, urlPath, params, nil
}
