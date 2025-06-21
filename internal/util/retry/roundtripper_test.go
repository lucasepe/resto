package retry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	callCount int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.callCount++

	// Simula una risposta JSON con success = false fino alla 3a chiamata
	success := m.callCount >= 3
	respBody := map[string]any{"success": success}
	bodyBytes, _ := json.Marshal(respBody)

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func TestRetryRoundTripper_JSONGval(t *testing.T) {

	retrier := NewRetrier(RetryOptions{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		MaxAttempts:  5,
	},
	)

	mock := &mockTransport{}

	rt := NewRoundTripperWithEval(mock, ".success == true", Exp(), retrier)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)

	resp, err := rt.RoundTrip(req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 3, mock.callCount) // ritenta fino a quando success == true
}

type mockTextTransport struct {
	callCount int
}

func (m *mockTextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.callCount++

	var body string
	if m.callCount < 4 {
		body = "please wait"
	} else {
		body = "done"
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
	}, nil
}
