# CS3103 G46 Go-Crawler

<!-- omit in toc -->
## Table of Contents

- [Components](#components)
  - [Crawler](#crawler)
    - [Usage](#usage)
- [Members](#members)
- [Acknowledgements](#acknowledgements)

## Components

### Crawler

A basic web crawler that crawls a given URL and returns a list of URLs found on the page.

#### Usage

> [!NOTE]
> Refer to [`example/crawler_report.json`](https://github.com/yusufaine/cs3103-gocrawler/blob/main/example/crawler_report.json) for an example of the output for the command below.

```bash
# Building the binary
make crawler

# Running the binary
./crawler --seed=https://example.com --depth=5 --report example/crawler_report.json
# ./crawler --help to see all options

# Without binary (requires Go 1.21+)
go run cmd/crawler/main.go --seed=https://example.com --depth=5 --report example/crawler_report.json
```

## Members

| **Name**              |
| --------------------- |
| Aryaa Adee Sandeep    |
| Jacob Kwan            |
| Ryan Aidan Jayasuriya |
| Yusuf Bin Musa        |

## Acknowledgements
