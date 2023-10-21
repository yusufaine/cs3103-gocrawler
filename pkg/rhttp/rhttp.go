// HTTP retry client with backoff and context support
package rhttp

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

const (
	defaultMaxRetryCount = 3
	defaultMinWaitMs     = 1000  // 1 second
	defaultMaxWaitMs     = 10000 // 10 seconds
	userAgent            = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
)

type Client struct {
	cl            *http.Client
	maxRetryCount int
	minWaitMs     int
	maxWaitMs     int
	retryPol      RetryPolicy
	backoffPol    BackoffPolicy
}

// By default, the client will retry 3 times with a linear backoff between 100ms
// and 1000ms. Refer to BackoffPolicy and RetryPolicy for more information.
func New(opts ...RHTTPOption) *Client {
	c := &Client{
		cl:            http.DefaultClient,
		maxRetryCount: defaultMaxRetryCount,
		minWaitMs:     defaultMinWaitMs,
		maxWaitMs:     defaultMaxWaitMs,
		retryPol:      DefaultRetry,
		backoffPol:    DefaultLinearBackoff,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// TODO: There may exist an issue if the request has a body
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", userAgent)
	for i := 0; i < c.maxRetryCount; i++ {
		select {
		case <-req.Context().Done():
			return resp, req.Context().Err()
		default:
			resp, err = c.cl.Do(req)
			retry, err := c.retryPol(resp, err)
			if !retry {
				return resp, err
			}
			wait := c.backoffPol(c.minWaitMs, c.maxWaitMs, i)
			log.Warn("retrying request", "attempt", i+1, "wait", wait, "link", req.URL.String())
			time.Sleep(wait)
		}
	}
	return resp, err
}
