package restclient

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

// dumpResponse copies the HTTP response body to the appropriate writer based on
// the response status code. If the status code indicates success (2xx), the body
// is copied to okWri. Otherwise, it is copied to koWri and an error is returned
// indicating the failure status.
//
// If copying the body fails, the function returns an error wrapping the cause.
//
// Both okWri and koWri must be non-nil writers.
//
// Example usage:
//
//	err := dumpResponse(resp, os.Stdout, os.Stderr)
func dumpResponse(res *http.Response, outwri, errwri io.Writer) error {
	if outwri == nil {
		outwri = io.Discard
	}

	if errwri == nil {
		errwri = io.Discard
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if res.Body == nil {
		if !statusOK {
			return fmt.Errorf("http request failed with status: %d %s", res.StatusCode, http.StatusText(res.StatusCode))
		}
		return nil
	}
	defer res.Body.Close()

	if !statusOK {
		_, err := io.Copy(errwri, res.Body)
		if err != nil {
			return fmt.Errorf("http status %d; also failed to read body: %w", res.StatusCode, err)
		}

		return fmt.Errorf("http request failed with status: %d %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	_, err := io.Copy(outwri, res.Body)

	return err
}

// isTextResponse returns true if the response appears to be a textual media type.
func isTextResponse(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return true
	}
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return strings.HasPrefix(media, "text/")
}
