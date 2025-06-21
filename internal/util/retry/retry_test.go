package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lucasepe/resto/internal/util/retry"
)

var (
	errStub error = errors.New("stub error")
)

func TestRetrierRetryContextDeadlineFail(t *testing.T) {
	r := retry.NewRetrier(
		retry.RetryOptions{
			InitialDelay: 125 * time.Millisecond,
			MaxDelay:     250 * time.Millisecond,
			MaxAttempts:  2,
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.Retry(ctx, retry.Exp(), func() (bool, error) {
		return true, nil
	})

	if err == nil {
		t.Fatal("unexpected nil error")
	}

	expectedErrorMessage := "context canceled"
	if err.Error() != expectedErrorMessage {
		t.Fatal(err)
	}
}

func TestRetrierRetry(t *testing.T) {
	r := retry.NewRetrier(
		retry.RetryOptions{
			InitialDelay: 125 * time.Millisecond,
			MaxDelay:     250 * time.Millisecond,
			MaxAttempts:  2,
		},
	)
	err := r.Retry(context.Background(), retry.Exp(), func() (bool, error) {
		return true, nil
	})

	if err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}
}

func TestRetrierRetryTriggerError(t *testing.T) {
	r := retry.NewRetrier(
		retry.RetryOptions{
			InitialDelay: 125 * time.Millisecond,
			MaxDelay:     250 * time.Millisecond,
			MaxAttempts:  2,
		},
	)
	err := r.Retry(context.Background(), retry.Exp(), func() (bool, error) {
		return false, errStub
	})

	if err == nil {
		t.Fatal("unexpected nil error")
	}

	if !errors.Is(err, errStub) {
		t.Fatal(err)
	}
}

func TestRetrierRetryFail(t *testing.T) {
	r := retry.NewRetrier(
		retry.RetryOptions{
			InitialDelay: 125 * time.Millisecond,
			MaxDelay:     250 * time.Millisecond,
			MaxAttempts:  2,
		},
	)

	err := r.Retry(context.Background(), retry.Exp(), func() (bool, error) {
		return false, nil
	})

	if err == nil {
		t.Fatal("unexpected nil error")
	}

	expectedErrorMessage := "function never succeeded in Retry"
	if err.Error() != expectedErrorMessage {
		t.Fatal(err)
	}
}
