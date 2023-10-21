package crawler

import (
	"flag"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/logger"
)

type Config struct {
	BlacklistHosts map[string]struct{}
	MaxDepth       int
	MaxRetries     int
	MaxRPS         float64
	RelReportPath  string
	SeedURL        string
	Timeout        time.Duration
}

func NewFlagConfig() *Config {
	var (
		c       Config
		blHosts string
		verbose bool
	)
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist, hosts will be blacklisted with and without 'www.' prefix")
	flag.IntVar(&c.MaxDepth, "depth", 20, "Max depth from seed")
	flag.StringVar(&c.SeedURL, "seed", "", "Seed URL, required (e.g https://example.com)")
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 15, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 5*time.Second, "Timeout for HTTP requests")
	flag.BoolVar(&verbose, "verbose", false, "For devs -- verbose logging, includes debug and short caller info")
	flag.StringVar(&c.RelReportPath, "report", "crawler_report.json", "Relative path to report file")
	flag.Parse()

	logger.Setup(verbose)
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

	c.MustValidate()

	log.Info("Running with config: ")
	log.Info("  ", "seed", c.SeedURL)
	log.Info("  ", "depth", c.MaxDepth)
	log.Info("  ", "blacklist", strings.Join(strings.Split(blHosts, ","), ", "))
	log.Info("  ", "retries", c.MaxRetries)
	log.Info("  ", "rps", c.MaxRPS)
	log.Info("  ", "timeout", c.Timeout)
	log.Info("  ", "verbose", verbose)
	log.Info("  ", "report", c.RelReportPath)

	return &c
}

func (c *Config) MustValidate() {
	if c.SeedURL == "" {
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
