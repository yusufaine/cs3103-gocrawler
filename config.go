package crawler

import (
	"flag"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/crawler/internal/logger"
)

// This file contains the necessary config for the crawler

type Config struct {
	BlacklistHosts map[string]struct{}
	MaxDepth       int
	MaxRetries     int
	MaxRPS         float64
	ProxyURL       *url.URL
	SeedURL        *url.URL
	Timeout        time.Duration
}

func SetupConfig() *Config {
	var (
		c       Config
		blHosts string
		proxy   string
		seed    string
		verbose bool
	)
	flag.IntVar(&c.MaxDepth, "depth", 10, "Max depth from seed")
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 15, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 5*time.Second, "Timeout for HTTP requests")
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist, hosts will be blacklisted with and without 'www.' prefix")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL (e.g http://localhost:8080)")
	flag.StringVar(&seed, "seed", "", "Seed URL, required (e.g https://example.com)")
	flag.BoolVar(&verbose, "verbose", false, "For devs -- verbose logging, includes debug and short caller info")
	flag.Parse()
	logger.Setup(verbose)

	c.SeedURL, _ = url.Parse(seed)
	c.ProxyURL, _ = url.Parse(proxy)

	c.BlacklistHosts = make(map[string]struct{})
	for _, host := range strings.Split(blHosts, ",") {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		c.BlacklistHosts[host] = struct{}{}

		if strings.HasPrefix(host, "www.") {
			c.BlacklistHosts[host[4:]] = struct{}{}
		} else {
			c.BlacklistHosts["www."+host] = struct{}{}
		}
	}

	return &c
}

func (c *Config) MustValidate() {
	if c.SeedURL.String() == "" {
		panic("--seed is required!")
	} else if c.MaxDepth < 1 {
		panic("--depth must be >= 1")
	} else if c.MaxDepth > 10 {
		log.Warn("'--depth' > 10 may take a long time to complete")
	} else if c.MaxRPS <= 0 {
		panic("--rps must be > 0")
	} else if c.MaxRPS > 20 {
		log.Warn("'--rps' > 20 may cause unexpected behaviour")
	} else if c.Timeout <= 0 {
		panic("--timeout must be > 0")
	} else if c.MaxRetries < 0 {
		panic("--retries must be >= 0")
	}
}

func (c *Config) PrintConfig() {
	blHosts := make([]string, 0, len(c.BlacklistHosts))
	for k := range c.BlacklistHosts {
		blHosts = append(blHosts, k)
	}
	slices.Sort(blHosts)

	log.Info("Running with config: ")
	log.Info(" ", "seed", c.SeedURL)
	log.Info(" ", "depth", c.MaxDepth)
	log.Info(" ", "proxy", c.ProxyURL)
	log.Info(" ", "blacklist", strings.Join(blHosts, ", "))
	log.Info(" ", "retries", c.MaxRetries)
	log.Info(" ", "rps", c.MaxRPS)
	log.Info(" ", "timeout", c.Timeout)
}
