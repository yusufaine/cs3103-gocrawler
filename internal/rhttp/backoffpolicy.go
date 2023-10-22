package rhttp

import (
	"math/rand"
	"sync"
	"time"
)

// make randomness deterministic
var rng = rand.New(rand.NewSource(3230))
var rngMutex sync.Mutex

func WithRNGSeed(seed int64) {
	rng = rand.New(rand.NewSource(seed))
}

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
	rngMutex.Lock()
	defer rngMutex.Unlock()
	waitMs += rng.Intn(minMs)
	return time.Duration(waitMs) * time.Millisecond
}

// ExponentialBackoff performs an exponential backoff based on the attempt number.
func ExponentialBackoff(minMs, maxMs, attempt int) time.Duration {
	waitMs := minMs * (1 << uint(attempt))
	if waitMs > maxMs {
		waitMs = maxMs
	}
	rngMutex.Lock()
	defer rngMutex.Unlock()
	waitMs += rng.Intn(minMs)
	return time.Duration(waitMs) * time.Millisecond
}
