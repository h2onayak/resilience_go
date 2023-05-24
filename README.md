# Resilience
For API documentation, refer to [GoDoc](https://pkg.go.dev/github.com/sony/gobreaker)

## What does it do?
1) Circuit Break and Fault Tolerance
2) Retrires

## Usage

### Create circuit break client with retry
```
	knowError := errors.New("error but sill can be consider as success")

	backoffInterval := 2 * time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	constantBackoff := resilience.NewConstantBackoff(backoffInterval, maximumJitterInterval)
	retrier := resilience.NewRetrier(constantBackoff)

	readyToTripOption := func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		// total number of rquest are greater than or equal to 10 and failure rate is greater than or equal to 60%, then trip 
		the curuite break to open state.
		return counts.Requests >= 10 && failureRatio >= 0.6
	}

	relianceClient := resilience.NewClient(
		resilience.WithCommandName("current-svc_upstream-svc"),
		resilience.WithHTTPTimeout(1*time.Second),
		resilience.WithHalfOpenStateMaxRequests(5),
		resilience.WithCloseStateInterval(10*time.Minute),
		resilience.WithOpenStateTimeout(1*time.Minute),
		resilience.WithReadyToTrip(readyToTripOption),
		resilience.WithOnStateChange(func(name string, from gobreaker.State, to gobreaker.State) {
			//log on client side for state change info on the requests
		}),
		resilience.WithIsSuccessful(func(err error) bool {
			if errors.Is(err, knowError) {
				return true
			}
			return false
		}),
		resilience.WithRetryCount(2),
		resilience.WithRetrier(retrier),
	)


	req, err := http.NewRequest("GET", "https://www.google.com/", nil)
	resp, err := relianceClient.Do(req)
```

## TODO
1. Add option to turn off circuit breaker client and just use http client (i.e WithEnableCircuitBreaker(default=true)).
2. Add plugin types to support different types metric collector (like, StasD, Newrelic, Prometheus)
3. Capture metrics for request latency, error and status codes.
4. Capture state change events and closed state interval metrics.
