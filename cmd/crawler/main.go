package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"

	"github.com/yusufaine/cs3203-g46-crawler/internal/crawler"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/fileexporter"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/logger"
)

func main() {
	// ensures that the data collected so far is exported when the program is terminated
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
		crawler.WithMaxRequestsPerSecond(config.MaxRPS),
	)
	cr.Crawl(ctx, config.SeedURL, 0)
	defer exportFiles(config, cr)
	log.Info("crawl completed")
}

func exportFiles(config *crawler.Config, cr *crawler.Crawler) {
	netinfo, err := json.MarshalIndent(cr.VisitedNetInfo, "", "  ")
	if err != nil {
		log.Error("unable to marshal network info", "error", err)
	}
	if err := fileexporter.WriteToFile(netinfo, config.RelNetInfoPath); err != nil {
		log.Error("unable to write to file", "file", config.RelNetInfoPath, "error", err)
	} else {
		log.Info("exported network info", "file", config.RelNetInfoPath)
	}

	pageinfo, err := json.MarshalIndent(cr.VisitedPageResp, "", "  ")
	if err != nil {
		log.Error("unable to marshal page info", "error", err)
	}
	if err := fileexporter.WriteToFile(pageinfo, config.RelPagesInfoPath); err != nil {
		log.Error("unable to write to file", "file", config.RelNetInfoPath, "error", err)
	} else {
		log.Info("exported page info", "file", config.RelPagesInfoPath)
	}
}
