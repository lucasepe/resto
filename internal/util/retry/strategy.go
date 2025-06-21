package retry

import (
	"math/rand"
	"time"
)

type Strategy interface {
	Policy(current, max time.Duration) time.Duration
}

func Exp() Strategy {
	return &expPolicy{curve: 2.0}
}

func Jittered(maxJitter time.Duration) Strategy {
	return &jitteredExp{
		curve:     2.0,
		maxJitter: maxJitter,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type expPolicy struct {
	curve float64
}

func (po *expPolicy) Policy(current, max time.Duration) time.Duration {
	next := time.Duration(float64(current) * po.curve)
	if next > max {
		return max
	}
	return next
}

type jitteredExp struct {
	curve     float64
	maxJitter time.Duration
	rng       *rand.Rand
}

func (je *jitteredExp) Policy(current, max time.Duration) time.Duration {
	expBackoff := time.Duration(float64(current) * je.curve)
	jitter := time.Duration(je.rng.Float64() * float64(je.maxJitter))
	next := expBackoff + jitter
	if next > max {
		return max
	}
	return next
}
