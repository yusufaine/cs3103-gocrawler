# gocrawler

> A simple concurrent webcrawler package written in Go.

## Components

| Component | Description                                                                                                |
| --------- | ---------------------------------------------------------------------------------------------------------- |
| `crawler` | Main crawler logic with a customisable `LinkExtractor` to allow users to determine how links are extracted |
| `logger`  | Sets up [`charmbracelet/log`](https://github.com/charmbracelet/log) to make logging less boring            |
| `rhttp`   | Wrapper over `net/http` with provided backoff and retry policies that can be customised                    |

## Usage

Examples of how to use the crawler package can be found in the `example` directory.
