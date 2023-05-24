package resilience_go

import (
	"time"

	"github.com/sony/gobreaker"
)

type Option func(*cbClient)

func WithCommandName(name string) Option {
	return func(c *cbClient) {
		c.settings.Name = name
	}
}

func WithHTTPTimeout(httpTimeout time.Duration) Option {
	return func(c *cbClient) {
		c.httpTimeout = httpTimeout
	}
}

func WithHalfOpenStateMaxRequests(maxRequests uint32) Option {
	return func(c *cbClient) {
		c.settings.MaxRequests = maxRequests
	}
}

func WithCloseStateInterval(interval time.Duration) Option {
	return func(c *cbClient) {
		c.settings.Interval = interval
	}
}

func WithOpenStateTimeout(timeout time.Duration) Option {
	return func(c *cbClient) {
		c.settings.Timeout = timeout
	}
}

func WithReadyToTrip(readyToTrip func(counts gobreaker.Counts) bool) Option {
	return func(c *cbClient) {
		c.settings.ReadyToTrip = readyToTrip
	}
}

func WithOnStateChange(onStateChange func(name string, from gobreaker.State, to gobreaker.State)) Option {
	return func(c *cbClient) {
		c.settings.OnStateChange = onStateChange
	}
}

func WithIsSuccessful(isSuccessful func(err error) bool) Option {
	return func(c *cbClient) {
		c.settings.IsSuccessful = isSuccessful
	}
}

func WithRetryCount(retryCount int) Option {
	return func(c *cbClient) {
		c.retryCount = retryCount
	}
}

func WithRetrier(retrier Retriable) Option {
	return func(c *cbClient) {
		c.retrier = retrier
	}
}

// TODO: Add metric collector to capture duration of each requests
//func WithMetricCollector(metricCollector MetricCollector) Option {
//	return func(c *cbClient) {
//		c.metricCollector = metricCollector
//	}
//}
