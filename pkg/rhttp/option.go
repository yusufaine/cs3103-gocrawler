package rhttp

type RHTTPOption func(*Client)

// WithBackoffPolicy sets the backoff policy for the client.
func WithBackoffPolicy(bp BackoffPolicy) RHTTPOption {
	return func(c *Client) {
		c.backoffPol = bp
	}
}

// WithRetryPolicy sets the retry policy for the client.
func WithRetryPolicy(rp RetryPolicy) RHTTPOption {
	return func(c *Client) {
		c.retryPol = rp
	}
}
