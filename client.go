package resilience_go

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	newrelic "github.com/newrelic/go-agent"
	"github.com/sony/gobreaker"
)

const (
	defaultRetryCount  = 0
	defaultHTTPTimeout = 1 * time.Second
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type cbClient struct {
	settings       *gobreaker.Settings
	circuitBreaker *gobreaker.CircuitBreaker
	httpClient     *http.Client
	httpTimeout    time.Duration
	retryCount     int
	retrier        Retriable
}

var err5xx = errors.New("server returned 5xx status code")

func NewClient(opts ...Option) Client {
	client := &cbClient{
		// Default ReadyToTrip returns true when the number of consecutive failures is more than 5.
		settings:    &gobreaker.Settings{},
		httpTimeout: defaultHTTPTimeout,
		retryCount:  defaultRetryCount,
		retrier:     NewNoRetrier(),
	}

	for _, opt := range opts {
		opt(client)
	}

	if client.httpClient == nil {
		client.httpClient = &http.Client{
			Timeout: client.httpTimeout,
		}
	}
	client.circuitBreaker = gobreaker.NewCircuitBreaker(*client.settings)

	return client
}

func (c *cbClient) Do(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	var bodyReader *bytes.Reader

	if request.Body != nil {
		reqData, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		// prevents closing the body between retries
		request.Body = ioutil.NopCloser(bodyReader)
	}

	for i := 0; i <= c.retryCount; i++ {
		if response != nil {
			response.Body.Close()
		}

		c.captureExternalSegment(request)

		cbResponse, err := c.circuitBreaker.Execute(func() (interface{}, error) {
			httpResp, err := c.httpClient.Do(request)
			if bodyReader != nil {
				_, err = bodyReader.Seek(0, 0)
				if err != nil {
					return nil, fmt.Errorf("failed to seek body: %v", err)
				}
			}
			if err != nil {
				return nil, err
			}
			defer httpResp.Body.Close()
			_, err = io.Copy(ioutil.Discard, httpResp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to copy response body: %v", err)
			}
			//5xx status code can be retried, lets return an error
			if httpResp.StatusCode >= http.StatusInternalServerError {
				return httpResp, err5xx
			}
			return httpResp, nil
		})
		if cbResponse != nil {
			response = cbResponse.(*http.Response)
		}
		if err != nil {
			backoffTime := c.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		break
	}

	//still there is 5xx error even after retry, lets return http response only.
	if err == err5xx {
		return response, nil
	}
	return response, err
}

func (c *cbClient) captureExternalSegment(request *http.Request) {
	txn := newrelic.FromContext(request.Context())
	if txn != nil {
		c.httpClient.Transport = newrelic.NewRoundTripper(txn, c.httpClient.Transport)
	}
}
