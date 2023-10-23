package gocrawler

import (
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// This file contains the necessary config for the crawler

type Config struct {
	BlacklistHosts map[string]struct{} // hosts to blacklist
	MaxDepth       int                 // max depth from seed
	MaxRetries     int                 // max retries for HTTP requests
	MaxRPS         float64             // max requests per second
	ProxyURL       *url.URL            // proxy URL, if any. useful to avoid IP bans
	SeedURLs       []*url.URL          // where to start crawling from
	Timeout        time.Duration       // timeout for HTTP requests
}

func (c *Config) MustValidate() {
	if len(c.SeedURLs) == 0 {
		panic("--seed is required!")
	}
	if c.MaxDepth < 1 {
		panic("--depth must be >= 1")
	}
	if c.MaxRPS <= 0 {
		panic("--rps must be > 0")
	}
	if c.Timeout <= 0 {
		panic("--timeout must be > 0")
	}
	if c.MaxRetries < 0 {
		panic("--retries must be >= 0")
	}
	if c.MaxDepth > 5 {
		log.Warn("depth is set to greater than 5 may take a long time to complete")
	}
	if c.MaxRPS > 20 {
		log.Warn("rps is set tp greater than 20 may cause unexpected behaviour such as rate limiting and IP bans")
	}
	if len(c.ProxyURL.String()) > 0 {
		log.Warn("proxy is set, this may affect the network info collected")
	}
}

func (c *Config) PrintConfig() {
	blHosts := make([]string, 0, len(c.BlacklistHosts))
	for k := range c.BlacklistHosts {
		blHosts = append(blHosts, k)
	}
	slices.Sort(blHosts)

	seeds := make([]string, 0, len(c.SeedURLs))
	for _, s := range c.SeedURLs {
		seeds = append(seeds, s.String())
	}
	slices.Sort(seeds)

	log.Info("Running with config: ")
	log.Info(" ", "seed", strings.Join(seeds, ", "))
	log.Info(" ", "depth", c.MaxDepth)
	log.Info(" ", "proxy", c.ProxyURL)
	log.Info(" ", "blacklist", strings.Join(blHosts, ", "))
	log.Info(" ", "retries", c.MaxRetries)
	log.Info(" ", "rps", c.MaxRPS)
	log.Info(" ", "timeout", c.Timeout)
}
