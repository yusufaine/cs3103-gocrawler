package rhttp

import (
	"net/http"
	"net/url"
	"time"
)

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

func WithProxy(proxyURL *url.URL) RHTTPOption {
	return func(c *Client) {
		if proxyURL.String() == "" {
			return
		}
		c.cl.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
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
