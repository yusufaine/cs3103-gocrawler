package sitemapper

import (
	"flag"
	"fmt"
	"math"
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

// SetupConfig wraps the gocrawler.Config and adds an additional report path field
// where users can specify where to export the report to.
func SetupConfig() *Config {
	var (
		c       Config
		blHosts string
		seed    string
		proxy   string
		verbose bool
	)
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 20, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 10*time.Second, "Timeout for HTTP requests")
	flag.StringVar(&c.ReportPath, "report", "", "Path to export report to. Defaults to 'sitemap_<seed>.json")
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist, hosts will be blacklisted with and without 'www.' prefix")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL")
	flag.StringVar(&seed, "seed", "", "Seed URL, required (e.g https://example.com)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging, includes short caller info")
	flag.Parse()
	logger.Setup(verbose)

	// sitemapper crawls indefinitely as long as the host is the same
	c.MaxDepth = math.MaxInt

	c.SeedURLs = strings.Split(seed, ",")

	// Parse proxy URL, if any
	parsedProxy, _ := url.Parse(proxy)
	c.ProxyURL = parsedProxy

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

	c.mustValidate()

	return &c
}

func (c *Config) mustValidate() {
	if len(c.SeedURLs) == 0 {
		panic("--seed is required!")
	}

	if len(c.SeedURLs) > 1 {
		panic("--seed expects only 1 URL, ensure that the value does not contain commas")
	}

	if parsedSeed, err := url.Parse(c.SeedURLs[0]); err != nil {
		panic("--seed is not valid!")
	} else if c.ReportPath == "" {
		c.ReportPath = fmt.Sprintf("%s_sitemap.json", parsedSeed.Host)
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

	log.Info("Running with config (ctrl-c to cancel crawling): ")
	log.Info(" ", "seed", strings.Join(c.SeedURLs, ", "))
	log.Info(" ", "proxy", c.ProxyURL)
	log.Info(" ", "retries", c.MaxRetries)
	log.Info(" ", "rps", c.MaxRPS)
	log.Info(" ", "timeout", c.Timeout)
	log.Info(" ", "report", c.ReportPath)
}
