# CS3103 Golang Webcrawler

This repo contains the source code for a generic parallel webcrawler written in Golang, which we then built upon to create a webcrawler to analyse the relevance of each country and region when it comes to the topic of "The International", a global DOTA 2 tournament, over the past few years based on what can be found on [Liquipedia](https://liquipedia.net/dota2/The_International).

<!-- omit in toc -->
## Table of Contents

- [Components](#components)
  - [Crawler](#crawler)
    - [Standalone Usage](#standalone-usage)
    - [Package Usage](#package-usage)
- [Members](#members)
- [Acknowledgements](#acknowledgements)

## Components

### Crawler

A basic web crawler that crawls a given URL and returns a list of URLs found on the page.

#### Standalone Usage

> [!NOTE]
> **Purpose**: Output results of crawling.
> Refer to [`example/crawler_report.json`](https://github.com/yusufaine/cs3103-gocrawler/blob/main/example/crawler_report.json) for an example of the output for the command below.

```bash
# Running the binary
./crawler \
  --seed=https://example.com/ \
  --depth=3 \
  --report=example/crawler_report.json \
  --rps=100 \
  --bl=pti.icann.org,facebook.com

# ./crawler --help to see all options

# Without binary (requires Go 1.21+)
go run cmd/crawler/main.go \
  --seed=https://example.com/ \
  --depth=3 \
  --report=example/crawler_report.json \
  --rps=100 \
  --bl=pti.icann.org,facebook.com

# Output logs to local file, useful for debugging
./crawler --seed=https://example.com --depth=3 2>&1 --v | tee crawler.log
```

#### Package Usage

> [!NOTE]
> **Purpose**: Repurpose the extracted information from crawler for other purposes.

```go
// Custom retry client that wraps over net/http's client
cr := crawler.New(
  ctx,
  config.MaxDepth,
  config,
  crawler.WithLinkExtractor(crawler.DefaultLinkExtractor),
  crawler.WithMaxRequestsPerSecond(config.MaxRPS),
)

// Start crawling
cr.Crawl(ctx, config.SeedURL, 0)

// Generate report if that's your intention
cr.GenerateReport(config)

// For consumption purposes
_ = cr.VisitedNetInfo
_ = cr.VisitedPageInfo
```

## Members

| **Name**              |
| --------------------- |
| Aryaa Adee Sandeep    |
| Jacob Kwan            |
| Ryan Aidan Jayasuriya |
| Yusuf Bin Musa        |

## Acknowledgements
