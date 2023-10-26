package tianalyser

import (
	"flag"
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

// Read config from flags to setup the crawler
func SetupConfig() *Config {
	var (
		c       Config
		proxy   string
		verbose bool
	)
	flag.IntVar(&c.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&c.MaxRPS, "rps", 0.5, "Max requests per second")
	flag.DurationVar(&c.Timeout, "timeout", 10*time.Second, "Timeout for HTTP requests")
	flag.StringVar(&c.ReportPath, "report", "ti_stats.json", "Path to export report to")
	flag.StringVar(&proxy, "proxy", "", "Proxy URL (e.g http://localhost:8080)")
	flag.BoolVar(&verbose, "verbose", false, "For devs -- verbose logging, includes debug and short caller info")
	flag.Parse()
	logger.Setup(verbose)

	c.MaxDepth = math.MaxInt

	// Parse proxy URL, if any
	c.ProxyURL, _ = url.Parse(proxy)

	c.mustValidate()

	return &c
}

func (c *Config) mustValidate() {
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

// Sanity check
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
