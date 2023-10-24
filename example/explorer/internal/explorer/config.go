package explorer

import (
	"flag"
	"fmt"
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

func SetupConfig() *Config {
	var (
		c       Config
		blHosts string
		seeds   string
		proxy   string
		verbose bool
	)

	// YYYY-MM-DD_HH-MM
	defaultReport := fmt.Sprintf("explorer_%s.json", time.Now().Format("2006-01-02_15-04"))

	flag.IntVar(&c.MaxDepth, "depth", 5, "Max depth from seed")
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 20, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 10*time.Second, "Timeout for HTTP requests")
	flag.StringVar(&c.ReportPath, "report", defaultReport, "Path to export report to")
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist, hosts will be blacklisted with and without 'www.' prefix")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL")
	flag.StringVar(&seeds, "seed", "", "Comma separated seed URL(s), required (e.g https://example.com)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging, includes short caller info")
	flag.Parse()
	logger.Setup(verbose)

	c.SeedURLs = strings.Split(seeds, ",")

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
	log.Info(" ", "depth", c.MaxDepth)
	log.Info(" ", "proxy", c.ProxyURL)
	log.Info(" ", "blacklist", strings.Join(blHosts, ", "))
	log.Info(" ", "retries", c.MaxRetries)
	log.Info(" ", "rps", c.MaxRPS)
	log.Info(" ", "timeout", c.Timeout)
	log.Info(" ", "report", c.ReportPath)
}
