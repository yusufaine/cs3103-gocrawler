# CS3103 Golang Webcrawler

This repo contains the source code for a generic parallel webcrawler written in Golang. In this directory there are 2 examples of how the crawler can be used as a package:

1. Sitemap generator, and
2. [Liquipedia](https://liquipedia.net/dota2/The_International) webcrawler

As part of our CS3103 mini-project, we continued building on top of the webcrawler to analyse the relevance of each country and region when it comes to the topic of "The International", a global DOTA 2 tournament, over the past few years based on what can be found on.

<!-- omit in toc -->
## Table of Contents

- [Components](#components)
  - [`crawler`](#crawler)
  - [`rhttp`](#rhttp)
    - [Package Usage](#package-usage)
    - [Standalone Usage](#standalone-usage)
- [Members](#members)
- [Acknowledgements](#acknowledgements)

## Components

### `crawler`

A concurrent web crawler that crawls a given URL and returns a list of URLs found on the page based on a default `LinkExtractor` method which users can override.

### `rhttp`

A simple wrapper over `net/http` that provides a few default backoff and retry policies that can also easily extend to a user's need.

#### Package Usage

> [!NOTE]
> **Purpose**: Repurpose the extracted information from crawler for other purposes (e.g. `sitemapgenerator`).

```go
config := tianalyser.SetupConfig()
config.MustValidate()
config.PrintConfig()
time.Sleep(3 * time.Second)

cr := gocrawler.New(ctx,
  &config.Config,
  []gocrawler.ResponseMatcher{gocrawler.HtmlContentFilter},
)

// Generates report using the collected data from the crawler
//  - cr.VisitedNetInfo
//  - cr.VisitedPageInfo
defer tianalyser.Generate(cr, config)
```

#### Standalone Usage

> [!NOTE]
> Refer to [`example/sitemapgenerator/sitemap.json`](https://github.com/yusufaine/cs3103-gocrawler/blob/main/example/crawler_report.json) for an example of the output for the command below.

```bash
# Running the binary
./tianalyser \
  --seed=https://liquipedia.net/dota2/The_International \ 
  --depth=5 \ 
  --report=example/tianalyser/ti_stats.json \ 
  --rps=10 

# ./tianalyser --help to see all options

# Without binary (requires Go 1.21+)
go run example/tianalyser/main.go \ 
  --seed=https://liquipedia.net/dota2/The_International \ 
  --depth=5 \ 
  --report=example/tianalyser/ti_stats.json \ 
  --rps=10 
```

## Members

| **Name**              |
| :-------------------- |
| Aryaa Adee Sandeep    |
| Jacob Kwan            |
| Ryan Aidan Jayasuriya |
| Yusuf Bin Musa        |

## Acknowledgements

- [Liquipedia](https://liquipedia.net/dota2/The_International) for providing the data for our project
- [Example usage of goquery](https://www.flysnow.org/2018/01/20/golang-goquery-examples-selector) for helping us understand how to use goquery
- [charmbracelet/log](https://github.com/charmbracelet/log) for making logging less boring
