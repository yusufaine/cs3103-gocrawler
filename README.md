# CS3103 G46 Go-Crawler

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
./crawler --seed=https://example.com --depth=3 --report example/crawler_report.json
# ./crawler --help to see all options

# Output logs to local file, useful for debugging
./crawler --seed=https://example.com --depth=3 2>&1 --v | tee crawler.log

# Without binary (requires Go 1.21+)
go run cmd/crawler/main.go --seed=https://example.com --depth=3 --report example/crawler_report.json
```

#### Package Usage

> [!NOTE]
> **Purpose**: Repurpose the extracted information from crawler for other purposes.

```go
// Custom retry client that wraps over net/http's client
retryClient := rhttp.New(
  rhttp.WithBackoffPolicy(rhttp.DefaultLinearBackoff),
  rhttp.WithMaxRetries(config.MaxRetries),
  rhttp.WithRetryPolicy(rhttp.DefaultRetry),
  rhttp.WithTimeout(config.Timeout),
)

// Create a new crawler
cr := crawler.New(
  config.MaxDepth,
  crawler.WithBlacklist(config.BlacklistHosts),
  crawler.WithLinkExtractor(crawler.DefaultLinkExtractor),
  crawler.WithMaxRequestsPerSecond(config.MaxRPS),
  crawler.WithRHttpClient(retryClient),
)

// Start crawling
cr.Crawl(ctx, config.SeedURL, 0)

// For consumption purposes
_ = cr.VisitedNetInfo
_ = cr.VisitedPageResp
```

## Members

| **Name**              |
| --------------------- |
| Aryaa Adee Sandeep    |
| Jacob Kwan            |
| Ryan Aidan Jayasuriya |
| Yusuf Bin Musa        |

## Acknowledgements
