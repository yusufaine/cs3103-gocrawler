package tianalyser

import (
	"flag"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler"
	"github.com/yusufaine/gocrawler/internal/logger"
)

type Config struct {
	gocrawler.Config
	ReportPath string
}

// Read config from flags to setup the crawler
func SetupConfig() *Config {
	var (
		c       Config
		blHosts string
		seeds   string
		proxy   string
		verbose bool
	)
	flag.IntVar(&c.MaxDepth, "depth", 5, "Max depth from seed")
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 20, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 10*time.Second, "Timeout for HTTP requests")
	flag.StringVar(&c.ReportPath, "report", "ti_stats.json", "Path to export report to")
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist, hosts will be blacklisted with and without 'www.' prefix")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL (e.g http://localhost:8080)")
	flag.StringVar(&seeds, "seed", "", "Comma separated seed URL(s), required (e.g https://example.com); invalid URLs will be ignored")
	flag.BoolVar(&verbose, "verbose", false, "For devs -- verbose logging, includes debug and short caller info")
	flag.Parse()
	logger.Setup(verbose)

	c.SeedURLs = strings.Split(seeds, ",")

	// Parse proxy URL, if any
	c.ProxyURL, _ = url.Parse(proxy)

	// Parse blacklist hosts into set for fast lookups
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

// Panics if config is invalid
func (c *Config) MustValidate() {
	c.Config.MustValidate()
	if c.ReportPath == "" {
		panic("--report is required")
	}
}

// Sanity check
func (c *Config) PrintConfig() {
	blHosts := make([]string, 0, len(c.BlacklistHosts))
	for k := range c.BlacklistHosts {
		blHosts = append(blHosts, k)
	}
	slices.Sort(blHosts)

	c.Config.PrintConfig()
	log.Info(" ", "report", c.ReportPath)
}
