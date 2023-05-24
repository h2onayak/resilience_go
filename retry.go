package resilience_go

import "time"

type Retriable interface {
	NextInterval(retry int) time.Duration
}

type RetriableFunc func(retry int) time.Duration

func NewRetrierFunc(f RetriableFunc) Retriable {
	return f
}

func (f RetriableFunc) NextInterval(retry int) time.Duration {
	return f(retry)
}

type noRetrier struct {
}

func NewNoRetrier() Retriable {
	return &noRetrier{}
}

func (r *noRetrier) NextInterval(retry int) time.Duration {
	return 0 * time.Millisecond
}

type retrier struct {
	backoff Backoff
}

func NewRetrier(backoff Backoff) Retriable {
	return &retrier{
		backoff: backoff,
	}
}

func (r *retrier) NextInterval(retry int) time.Duration {
	return r.backoff.Next(retry)
}
