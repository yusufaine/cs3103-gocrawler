package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/crawler"
	"github.com/yusufaine/crawler/example/liquipediacrawler/internal/linkextractor"
	"github.com/yusufaine/crawler/example/sitemapgenerator/sitemap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	config := crawler.SetupConfig()
	config.MustValidate()
	parsedURL, err := url.Parse("https://liquipedia.net/dota2/The_International")
	if err != nil {
		panic(err)
	}
	config.SeedURL = parsedURL
	config.PrintConfig()
	time.Sleep(3 * time.Second)

	cr := crawler.New(ctx, config,
		[]crawler.ResponseMatcher{crawler.HtmlContentFilter})

	// TODO: replace this with a the proper liquipedia analytics generator
	defer sitemap.Generate(&sitemap.Config{
		Config:     *config,
		ReportPath: "sitemap.json",
	}, cr)

	go func() {
		defer func() {
			cancel()
			close(sig)
		}()

		<-ctx.Done()
		fmt.Println()
		log.Info("stopping crawler", "signal", <-sig)
	}()

	// TODO: use linkextractor.TIAnalyserLinkExtractor to exclude visited outgoing links
	cr.Crawl(ctx, linkextractor.ReportLinkExtractor, config.SeedURL, 0)
	log.Info("crawl completed")
}
