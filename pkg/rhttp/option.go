package rhttp

import "time"

type RHTTPOption func(*Client)

// WithBackoffPolicy sets the backoff policy for the client.
func WithBackoffPolicy(bp BackoffPolicy) RHTTPOption {
	return func(c *Client) {
		c.backoffPol = bp
	}
}

func WithMaxRetries(maxRetries int) RHTTPOption {
	return func(c *Client) {
		c.maxRetryCount = maxRetries
	}
}

// WithRetryPolicy sets the retry policy for the client.
func WithRetryPolicy(rp RetryPolicy) RHTTPOption {
	return func(c *Client) {
		c.retryPol = rp
	}
}

func WithTimeout(timeout time.Duration) RHTTPOption {
	return func(c *Client) {
		c.cl.Timeout = timeout
	}
}
