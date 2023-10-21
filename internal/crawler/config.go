package crawler

import (
	"flag"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Config struct {
	BlacklistHosts map[string]struct{}
	MaxDepth       int
	SeedURL        string
	MaxRetries     int
	MaxRPS         float64
	Timeout        time.Duration
	Verbose        bool
	RelReportPath  string
}

func NewFlagConfig() *Config {
	var (
		config  Config
		blHosts string
	)
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist")
	flag.IntVar(&config.MaxDepth, "depth", 20, "Max depth from seed")
	flag.StringVar(&config.SeedURL, "seed", "", "Seed URL, required (e.g https://example.com)")
	flag.IntVar(&config.MaxRetries, "retries", 3, "Max retries for HTTP requests")
	flag.Float64Var(&config.MaxRPS, "rps", 15, "Max requests per second")
	flag.DurationVar(&config.Timeout, "timeout", 5*time.Second, "Timeout for HTTP requests")
	flag.BoolVar(&config.Verbose, "verbose", false, "For devs -- verbose logging, includes debug and short caller info")
	flag.StringVar(&config.RelReportPath, "report", "crawler_report.json", "Relative path to report file")
	flag.Parse()

	config.BlacklistHosts = make(map[string]struct{})
	for _, host := range strings.Split(blHosts, ",") {
		config.BlacklistHosts[host] = struct{}{}
	}

	return &config
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
	}
}

func (c *Config) PrintRunningConfig() {
	blhosts := make([]string, 0, len(c.BlacklistHosts))
	for k := range c.BlacklistHosts {
		blhosts = append(blhosts, k)
	}

	log.Info("Running with config: ")
	log.Info("  ", "seed", c.SeedURL)
	log.Info("  ", "depth", c.MaxDepth)
	log.Info("  ", "blacklist", blhosts)
	log.Info("  ", "retries", c.MaxRetries)
	log.Info("  ", "rps", c.MaxRPS)
	log.Info("  ", "timeout", c.Timeout)
	log.Info("  ", "verbose", c.Verbose)
	log.Info("  ", "report", c.RelReportPath)
}
