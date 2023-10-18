package keywordcounter

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	BlacklistHosts map[string]struct{}
	Keywords       []string
	MaxDepth       int
	OutputFile     string
	SeedURL        string
	Verbose        bool
}

func NewFlagConfig() *Config {
	var (
		blHosts    string
		config     Config
		csKeywords string
		outfile    string
	)
	outfile = fmt.Sprintf("keywords_%s.json", time.Now().Format(time.RFC3339))
	flag.StringVar(&blHosts, "blacklist", "", "Comma separated list of hosts to blacklist")
	flag.StringVar(&csKeywords, "keywords", "", "Comma separated list of keywords")
	flag.IntVar(&config.MaxDepth, "hops", 20, "Max hops from seed")
	flag.StringVar(&config.OutputFile, "output", outfile, "Output file name")
	flag.StringVar(&config.SeedURL, "seed", "", "Seed URL, required")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose logging, includes debug and caller info")
	flag.Parse()

	config.BlacklistHosts = make(map[string]struct{})
	for _, host := range strings.Split(blHosts, ",") {
		config.BlacklistHosts[host] = struct{}{}
	}
	config.Keywords = strings.Split(csKeywords, ",")

	return &config
}

func (c *Config) MustValidate() {
	if c.SeedURL == "" {
		panic("--seed is required!")
	} else if len(c.Keywords) == 0 {
		panic("--keywords is required!")
	}
}
