package crawler

import (
	"golang.org/x/time/rate"
)

type CrawlerOption func(*Crawler)

// Determines how links are extracted from the response body.
// See: linkextractor.go
func WithLinkExtractor(le LinkExtractor) CrawlerOption {
	return func(c *Crawler) {
		c.le = le
	}
}

// Determines the maximum number of requests per second for the crawler.
//
// Warning: Setting a value too high may cause the application to crash. This number is
// dependent on the host machine.
func WithMaxRequestsPerSecond(lim float64) CrawlerOption {
	return func(c *Crawler) {
		c.rl = rate.NewLimiter(rate.Limit(lim), 1)
	}
}
