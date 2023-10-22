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

<!-- TODO: Update this part to show the liquipediacrawler example -->
#### Package Usage

> [!NOTE]
> **Purpose**: Repurpose the extracted information from crawler for other purposes (e.g. `sitemapgenerator`).

```go
// sitemap.Config embeds crawler.Config
config := sitemap.SetupConfig()
config.MustValidate()

// Pass in crawler.Config
cr := crawler.New(ctx, &config.Config, config.MaxRPS)

// cr contains the extracted information from crawler
//  _ = cr.VisitedNetInfo   // network-related information
//  _ = cr.VisitedPageInfo  // page-related information
defer sitemap.Generate(config, cr)
```

<!-- TODO: Update this part to show the liquipediacrawler example -->
#### Standalone Usage

> [!NOTE]
> Refer to [`example/sitemapgenerator/sitemap.json`](https://github.com/yusufaine/cs3103-gocrawler/blob/main/example/crawler_report.json) for an example of the output for the command below.

```bash
# Running the binary
./sitemapgenerator \
  --seed=https://example.com/ \
  --depth=3 \
  --report=example/sitemapgenerator/sitemap.json \
  --rps=100 \
  --bl=pti.icann.org,facebook.com

# ./sitemapgenerator --help to see all options

# Without binary (requires Go 1.21+)
go run example/sitemapgenerator/main.go \
  --seed=https://example.com/ \
  --depth=3 \
  --report=example/sitemapgenerator/sitemap.json \
  --rps=100 \
  --bl=pti.icann.org,facebook.com
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
