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
	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
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

	retryClient := rhttp.New(
		rhttp.WithBackoffPolicy(rhttp.DefaultLinearBackoff),
		rhttp.WithMaxRetries(config.MaxRetries),
		rhttp.WithRetryPolicy(rhttp.DefaultRetry),
		rhttp.WithTimeout(config.Timeout),
	)

	cr := crawler.New(
		config.MaxDepth,
		crawler.WithBlacklist(config.BlacklistHosts),
		crawler.WithLinkExtractor(crawler.DefaultLinkExtractor),
		crawler.WithMaxRequestsPerSecond(config.MaxRPS),
		crawler.WithRHttpClient(retryClient),
	)

	cr.Crawl(ctx, config.SeedURL, 0)
	log.Info("crawl completed")
	exportReport(config, cr)
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
		VisitedPageResp map[string]crawler.PageInfo      `json:"page_info"`
	}{
		Seed:            config.SeedURL,
		Depth:           config.MaxDepth,
		Blacklist:       bls,
		VisitedNetInfo:  cr.VisitedNetInfo,
		VisitedPageResp: cr.VisitedPageResp,
	}
	for k, v := range cr.VisitedNetInfo {
		for i, v1 := range v {
			slices.Sort(v1.DNSAddrs)
			slices.Sort(v1.RemoteAddrs)
			slices.Sort(v1.VisitedPaths)
			filecontent.VisitedNetInfo[k][i] = v1
		}
	}

	if err := filewriter.ToJSON(filecontent, config.RelReportPath); err != nil {
		log.Error("unable to write to file",
			"file", config.RelReportPath,
			"error", err)
	} else {
		log.Info("exported crawler report", "file", config.RelReportPath)
	}
}
