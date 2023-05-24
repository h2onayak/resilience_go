package resilience_go

import (
	"math"
	"math/rand"
	"time"
)

type Backoff interface {
	Next(retry int) time.Duration
}

type constantBackoff struct {
	backoffInterval       int64
	maximumJitterInterval int64
}

func NewConstantBackoff(backoffInterval, maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}

	return &constantBackoff{
		backoffInterval:       int64(backoffInterval / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

func (cb *constantBackoff) Next(retry int) time.Duration {
	return (time.Duration(cb.backoffInterval) * time.Millisecond) + (time.Duration(rand.Int63n(cb.maximumJitterInterval+1)) * time.Millisecond)
}

type exponentialBackoff struct {
	exponentFactor        float64
	initialTimeout        float64
	maxTimeout            float64
	maximumJitterInterval int64
}

func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration, exponentFactor float64,
	maximumJitterInterval time.Duration) Backoff {
	// protect against panic when generating random jitter
	if maximumJitterInterval < 0 {
		maximumJitterInterval = 0
	}

	return &exponentialBackoff{
		exponentFactor:        exponentFactor,
		initialTimeout:        float64(initialTimeout / time.Millisecond),
		maxTimeout:            float64(maxTimeout / time.Millisecond),
		maximumJitterInterval: int64(maximumJitterInterval / time.Millisecond),
	}
}

func (eb *exponentialBackoff) Next(retry int) time.Duration {
	if retry < 0 {
		retry = 0
	}
	return time.Duration(math.Min(eb.initialTimeout*math.Pow(eb.exponentFactor, float64(retry)),
		eb.maxTimeout)+float64(rand.Int63n(eb.maximumJitterInterval+1))) * time.Millisecond
}
