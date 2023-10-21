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
	"github.com/yusufaine/cs3203-g46-crawler/internal/crawler"
	"github.com/yusufaine/cs3203-g46-crawler/internal/liquipediacrawler"
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

	cr := crawler.New(ctx, config, config.MaxRPS)
	defer cr.GenerateReport(config)

	go func() {
		defer func() {
			cancel()
			close(sig)
		}()

		<-ctx.Done()
		fmt.Println()
		log.Info("stopping crawler", "signal", <-sig)
	}()

	// TODO: use liquipediacrawler.TIAnalyserLinkExtractor to exclude visited outgoing links
	cr.Crawl(ctx, liquipediacrawler.ReportLinkExtractor, config.SeedURL, 0)
	log.Info("crawl completed")
}
