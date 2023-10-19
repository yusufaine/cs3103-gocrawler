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

```bash
# Building the binary
make crawler

# Running the binary
./crawler --seed=https://example.com --depth=20
# ./crawler --help to see all options

# Without binary (requires Go 1.21+)
go run cmd/crawler/main.go --seed=https://example.com --depth=20
```

## Members

| **Name**              |
| --------------------- |
| Aryaa Adee Sandeep    |
| Jacob Kwan            |
| Ryan Aidan Jayasuriya |
| Yusuf Bin Musa        |

## Acknowledgements
