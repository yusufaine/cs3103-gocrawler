package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler"
	"github.com/yusufaine/gocrawler/example/sitemapper/internal/sitemapper"
)

func main() {
	// ensures that the data collected so far is exported when the user terminates the program
	// (e.g. ctrl+c)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	// sitemapper.Config embeds gocrawler.Config
	config := sitemapper.SetupConfig()
	config.PrintConfig()
	time.Sleep(3 * time.Second)
	start := time.Now()

	cr := gocrawler.New(ctx,
		&config.Config,
		[]gocrawler.ResponseMatcher{gocrawler.IsHtmlContent})
	defer func() {
		log.Info("generating sitemap", "file", config.ReportPath)
		sitemapper.Generate(config, cr, time.Since(start))
	}()

	go func() {
		defer func() {
			cancel()
			close(sig)
		}()

		<-ctx.Done()
		fmt.Println()
		log.Info("stopping crawler, press ctrl+c again to force quit", "signal", <-sig)
	}()

	var wg sync.WaitGroup
	for _, seed := range config.SeedURLs {
		wg.Add(1)
		go func(seed string) {
			defer wg.Done()
			cr.Crawl(ctx, sitemapper.SameHostLinkExtractor, 0, seed, "")
		}(seed)
	}
	wg.Wait()
	log.Info("crawl completed")
}
