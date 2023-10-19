package rhttp_test

import (
	"testing"
	"time"

	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
)

// Testing RNG, do not run individually, use `make ci`
func TestDefaultLinearBackoff(t *testing.T) {
	min, max := 100, 1000
	exp := []time.Duration{134, 285, 361, 408, 544, 698, 715, 887, 903, 1031}
	for i, e := range exp {
		exp[i] = e * time.Millisecond
	}

	for i := 0; i < max/min; i++ {
		b := rhttp.DefaultLinearBackoff(min, max, i)
		if b != exp[i] {
			t.Errorf("Expected %d, got %d", exp[i], b)
		}
	}
}

// Testing RNG, do not run individually, use `make ci`
func TestExponentialBackoff(t *testing.T) {
	min, max := 100, 1_000_000
	// notice the doubling + jitter
	exp := []time.Duration{117, 216, 452, 893, 1684, 3212, 6417, 12872, 25627, 51214}
	for i, e := range exp {
		exp[i] = e * time.Millisecond
	}
	for i := 0; i < 10; i++ {
		b := rhttp.ExponentialBackoff(min, max, i)
		if b != exp[i] {
			t.Errorf("Expected %d, got %d", exp[i], b)
		}
	}
}
