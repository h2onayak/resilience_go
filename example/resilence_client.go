package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	resilience "github.com/dunzoit/dunzo_commons/resilience/v2"
	"github.com/sony/gobreaker"
)

func createClient(commandName string) resilience.Client {
	knowError := errors.New("known error to consider as success")

	backoffInterval := 2 * time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	constantBackoff := resilience.NewConstantBackoff(backoffInterval, maximumJitterInterval)
	retrier := resilience.NewRetrier(constantBackoff)

	readyToTripOption := func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 5 && failureRatio >= 0.6
	}

	return resilience.NewClient(
		resilience.WithCommandName(commandName),
		resilience.WithHTTPTimeout(1*time.Second),
		resilience.WithHalfOpenStateMaxRequests(5),
		resilience.WithCloseStateInterval(10*time.Minute),
		resilience.WithOpenStateTimeout(1*time.Minute),
		resilience.WithReadyToTrip(readyToTripOption),
		resilience.WithOnStateChange(func(name string, from gobreaker.State, to gobreaker.State) {
			//log on client side for state change info on requests
		}), resilience.WithIsSuccessful(func(err error) bool {
			if errors.Is(err, knowError) {
				return true
			}
			return false
		}),
		resilience.WithRetryCount(2), resilience.WithRetrier(retrier))
}

func main() {
	client := createClient("current-svc_upstream-svc")

	req, err := http.NewRequest("GET", "https://www.google.com/", nil)
	if err != nil {
		fmt.Printf("Error on creating http req: %v", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error on client connection: %v", err)
		return
	}

	fmt.Printf("resp: %v", resp)
}
