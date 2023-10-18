// HTTP retry client with backoff and context support
package rhttp

import (
	"net/http"
	"time"
)

const (
	defaultRetryCount = 3
	defaultMinWaitMs  = 100
	defaultMaxWaitMs  = 1000
	userAgent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
)

type Client struct {
	cl         *http.Client
	retryCount int
	minWaitMs  int
	maxWaitMs  int
	retryPol   RetryPolicy
	backoffPol BackoffPolicy
}

// By default, the client will retry 3 times with a linear backoff between 100ms
// and 1000ms. Refer to BackoffPolicy and RetryPolicy for more information.
func New(opts ...RHTTPOption) *Client {
	c := &Client{
		cl:         http.DefaultClient,
		retryCount: defaultRetryCount,
		minWaitMs:  defaultMinWaitMs,
		maxWaitMs:  defaultMaxWaitMs,
		retryPol:   DefaultRetry,
		backoffPol: DefaultLinearBackoff,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// TODO: There may exist an issue if the request has a body
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", userAgent)
	var resp *http.Response
	var err error
	for i := 0; i < c.retryCount; i++ {
		resp, err = c.cl.Do(req)
		retry, err := c.retryPol(resp, err)
		if !retry {
			return resp, err
		}
		wait := c.backoffPol(c.minWaitMs, c.maxWaitMs, i)
		time.Sleep(wait)
	}
	return resp, err
}
