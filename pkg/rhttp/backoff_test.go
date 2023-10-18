package rhttp_test

import (
	"testing"
	"time"

	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
)

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

func TestExponentialBackoff(t *testing.T) {
	min, max := 100, 1_000_000
	// notice the doubling + jitter
	exp := []time.Duration{134, 285, 461, 808, 1644, 3298, 6415, 12887, 25603, 51231}
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
