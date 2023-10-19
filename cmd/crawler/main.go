package main

import (
	"context"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/internal/crawler"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/filewriter"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/logger"
)

func main() {
	// ensures that the data collected so far is exported when the user terminates the program
	// (e.g. ctrl+c)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	config := crawler.NewFlagConfig()
	logger.Setup(config.Verbose)
	config.MustValidate()
	config.PrintRunningConfig()
	time.Sleep(3 * time.Second)

	cr := crawler.New(
		config.MaxDepth,
		crawler.WithBlacklist(config.BlacklistHosts),
		crawler.WithLinkExtractor(crawler.DefaultLinkExtractor),
		crawler.WithMaxRequestsPerSecond(config.MaxRPS),
	)
	cr.Crawl(ctx, config.SeedURL, 0)
	defer exportReport(config, cr)
	log.Info("crawl completed")
}

func exportReport(config *crawler.Config, cr *crawler.Crawler) {
	bls := make([]string, 0, len(config.BlacklistHosts))
	for k := range config.BlacklistHosts {
		bls = append(bls, k)
	}
	slices.Sort(bls)

	var filecontent = struct {
		Seed      string   `json:"seed"`
		Depth     int      `json:"max_depth"`
		Blacklist []string `json:"blacklist"`

		VisitedNetInfo  map[string][]crawler.NetworkInfo `json:"network_info"`
		VisitedPageResp map[string][]crawler.PageInfo    `json:"page_info"`
	}{
		Seed:            config.SeedURL,
		Depth:           config.MaxDepth,
		Blacklist:       bls,
		VisitedNetInfo:  cr.VisitedNetInfo,
		VisitedPageResp: cr.VisitedPageResp,
	}
	if err := filewriter.ToJSON(filecontent, config.RelReportPath); err != nil {
		log.Error("unable to write to file",
			"file", config.RelReportPath,
			"error", err)
	} else {
		log.Info("exported crawler info", "file", config.RelReportPath)
	}
}
