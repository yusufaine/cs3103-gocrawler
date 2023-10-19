package rhttp

import (
	"math/rand"
	"time"
)

// make randomness deterministic
var seed = rand.New(rand.NewSource(3230))

// Determines how long to wait before retrying a request bounded by the minimum
// and maximum wait times in milliseconds. Jitter should be applied to prevent
// thundering herds.
type BackoffPolicy func(minMs, maxMs, attempt int) time.Duration

// DefaultLinearBackoff performs a linear backoff based on the attempt number.
func DefaultLinearBackoff(minMs, maxMs, attempt int) time.Duration {
	waitMs := minMs + (attempt * 100)
	if waitMs > maxMs {
		waitMs = maxMs
	}
	waitMs += seed.Intn(minMs)
	return time.Duration(waitMs) * time.Millisecond
}

// ExponentialBackoff performs an exponential backoff based on the attempt number.
func ExponentialBackoff(minMs, maxMs, attempt int) time.Duration {
	waitMs := minMs * (1 << uint(attempt))
	if waitMs > maxMs {
		waitMs = maxMs
	}
	waitMs += seed.Intn(minMs)
	return time.Duration(waitMs) * time.Millisecond
}
