package retry

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/lucasepe/resto/internal/util/jq"
)

func NewRoundTripperWithEval(next http.RoundTripper, expr string, strategy Strategy, retrier Retrier) *retryRoundTripper {
	return &retryRoundTripper{
		retrier:    retrier,
		strategy:   strategy,
		next:       next,
		expression: expr,
	}
}

type retryRoundTripper struct {
	retrier    Retrier
	strategy   Strategy
	next       http.RoundTripper
	expression string
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	err := rt.retrier.Retry(context.Background(), rt.strategy, func() (bool, error) {
		var err error
		resp, err = rt.next.RoundTrip(req)
		if err != nil {
			return false, err
		}

		if resp.Body == nil {
			return false, nil
		}

		bin, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(bin)) // ripristina il body

		contentType := strings.ToLower(resp.Header.Get("Content-Type"))

		isJSON := strings.Contains(contentType, "application/json")

		switch {
		case rt.expression != "" && isJSON:
			return jq.EvalBoolExpr(bin, rt.expression)

		default:
			// Non gestito: consideriamo valido
			return true, nil
		}
	})

	return resp, err
}
