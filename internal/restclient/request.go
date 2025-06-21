package restclient

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// cloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func cloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = cloneHeader(req.Header)

	return r
}

// cloneHeader creates a deep copy of an http.Header.
func cloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}

// setHeaders sets HTTP headers on the given request from a list of "Key:Value" formatted strings.
//
// Parameters:
//   - req: The HTTP request to which headers will be added.
//   - headers: A variadic list of strings in the format "Key:Value". Leading/trailing spaces around keys and values are trimmed.
//     Headers with empty keys or invalid formats (no colon) are ignored with a warning logged to stderr.
//     Headers with empty values (e.g., "X-Debug:") are accepted and added.
//
// Example:
//
//	setHeaders(req, "Authorization: Bearer token", "X-Debug:")
func setHeaders(req *http.Request, headers ...string) {
	for _, el := range headers {
		idx := strings.Index(el, ":")
		if idx < 0 {
			fmt.Fprintf(os.Stderr, "ignoring invalid header format: %q\n", el)
			continue
		}

		key := strings.TrimSpace(el[:idx])
		val := strings.TrimSpace(el[idx+1:])

		if key == "" {
			fmt.Fprintf(os.Stderr, "ignoring header with empty key: %q\n", el)
			continue
		}

		req.Header.Set(key, val)
	}
}

// composeURL constructs a complete URL by combining a base URL, a relative path, and optional query parameters.
//
// Parameters:
//   - baseURL: The base URL, e.g., "https://example.com".
//   - path: The relative path to append, e.g., "api/v1/resource".
//   - params: An optional list of query parameters in the form "key:value".
//     Values can be empty (e.g., "debug:") to represent query strings like ?debug=.
//
// Returns:
//   - A string representing the full constructed URL.
//   - An error if the final URL cannot be parsed.
func composeURL(baseURL, path string, params ...string) (string, error) {
	uri := strings.TrimSuffix(baseURL, "/")
	if len(path) > 0 {
		uri = fmt.Sprintf("%s/%s", uri, strings.TrimPrefix(path, "/"))
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	query := u.Query()

	for _, el := range params {
		idx := strings.Index(el, ":")
		if idx <= 0 {
			// Query param non valido: ignora o logga un warning
			fmt.Fprintf(os.Stderr, "ignoring invalid query param format: %q\n", el)
			continue
		}

		key := strings.TrimSpace(el[:idx])
		val := strings.TrimSpace(el[idx+1:])

		if key == "" {
			fmt.Fprintf(os.Stderr, "ignoring query param with empty key: %q\n", el)
			continue
		}

		query.Add(key, val)
	}

	u.RawQuery = query.Encode()

	return u.String(), nil
}
