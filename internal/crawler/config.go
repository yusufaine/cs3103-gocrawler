package crawler

import (
	"flag"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Config struct {
	BlacklistHosts   map[string]struct{}
	MaxDepth         int
	SeedURL          string
	MaxRPS           float64
	Timeout          time.Duration
	Verbose          bool
	RelNetInfoPath   string
	RelPagesInfoPath string
}

func NewFlagConfig() *Config {
	var (
		config  Config
		blHosts string
	)
	flag.StringVar(&blHosts, "bl", "", "Comma separated list of hosts to blacklist")
	flag.IntVar(&config.MaxDepth, "depth", 20, "Max depth from seed")
	flag.StringVar(&config.SeedURL, "seed", "", "Seed URL, required (e.g https://example.com)")
	flag.Float64Var(&config.MaxRPS, "rps", 20, "Max requests per second")
	flag.DurationVar(&config.Timeout, "timeout", 5*time.Second, "Timeout for HTTP requests")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose logging, includes debug and caller info")
	flag.StringVar(&config.RelNetInfoPath, "out-net", "network_info.json", "Relative path to network info file")
	flag.StringVar(&config.RelPagesInfoPath, "out-page", "pages_info.json", "Relative path to page info file")
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
		log.Warn("--depth > 10 may take a long time to complete")
	} else if c.MaxRPS <= 0 {
		panic("--rps must be > 0")
	} else if c.MaxRPS > 20 {
		log.Warn("--rps > 20 may cause rate limiting")
	}
}

func (c *Config) PrintRunningConfig() {
	blhosts := make([]string, 0, len(c.BlacklistHosts))
	for k := range c.BlacklistHosts {
		blhosts = append(blhosts, k)
	}

	log.Info("Running with config",
		"blacklist", blhosts,
		"max_depth", c.MaxDepth,
		"seed", c.SeedURL,
		"timeout", c.Timeout,
		"verbose", c.Verbose,
		"net_info_rel_path", c.RelNetInfoPath,
		"page_info_rel_path", c.RelPagesInfoPath)
}
