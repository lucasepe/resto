package retry

import (
	"context"
	"errors"
	"time"

	"github.com/lucasepe/x/env"
)

var (
	ErrExhausted = errors.New("function never succeeded in Retry")
)

type Retrier interface {
	Retry(context.Context, Strategy, RetryFunc) error
}

type RetryFunc func() (bool, error)

type RetryOptions struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	MaxAttempts  int
	MaxJitter    time.Duration
}

func OptionsFromEnv() (res RetryOptions) {
	res.InitialDelay = env.Duration(initialDelayEnv, 500*time.Millisecond)
	res.MaxDelay = env.Duration(maxDelayEnv, 20*time.Second)
	res.MaxAttempts = env.Int(maxAttemptsEnv, 15)
	res.MaxJitter = env.Duration(maxJitterEnv, 1*time.Second)
	return res
}

func NewRetrier(opts RetryOptions) Retrier {
	ri := &retrierImpl{
		initialDelay: opts.InitialDelay,
		maxDelay:     opts.MaxDelay,
		maxAttempts:  opts.MaxAttempts,
	}

	// Sanity check: initialDelay should not be greater than maxDelay
	if ri.initialDelay > ri.maxDelay {
		ri.initialDelay = time.Duration(0.1 * float64(ri.maxDelay))
	}

	if ri.maxAttempts <= 0 {
		ri.maxAttempts = 10
	}

	return ri
}

const (
	initialDelayEnv = "INITIAL_DELAY"
	maxDelayEnv     = "MAX_DELAY"
	maxAttemptsEnv  = "MAX_ATTEMPTS"
	maxJitterEnv    = "MAX_JITTER"
)

type retrierImpl struct {
	initialDelay time.Duration
	maxDelay     time.Duration
	maxAttempts  int
}

func (ri *retrierImpl) Retry(ctx context.Context, strategy Strategy, fn RetryFunc) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		err  error
		done bool
	)

	interval := ri.initialDelay

	for i := 0; !done && i < ri.maxAttempts; i++ {
		//log.Printf("retry: attempt %d of %d\n", i+1, ri.attempts)

		done, err = fn()

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err != nil {
			return err
		}

		if !done && i+1 < ri.maxAttempts { // do not sleep after last attempt
			select {
			case <-time.After(interval):
				// continue
			case <-ctx.Done():
				return ctx.Err()
			}
			interval = strategy.Policy(interval, ri.maxDelay)
		}
	}

	if !done {
		return ErrExhausted
	}

	return nil
}
