package gocrawler

import (
	"net/url"
	"time"
)

// This file contains the necessary config for the crawler

type Config struct {
	BlacklistHosts map[string]struct{} // hosts to blacklist
	MaxDepth       int                 // max depth from seed
	MaxRetries     int                 // max retries for HTTP requests
	MaxRPS         float64             // max requests per second
	ProxyURL       *url.URL            // proxy URL, if any. useful to avoid IP bans
	SeedURLs       []string            // where to start crawling from
	Timeout        time.Duration       // timeout for HTTP requests
}
