# gocrawler

> A simple concurrent webcrawler package written in Go.

## Packages

| Package             | Description                                                                                                                                                         |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `gocrawler` (main)  | Main crawler logic with a customisable `LinkExtractor` to allow users to determine how links are extracted, and `ResponseMatcher` to filter out unwanted responses. |
| `logger` (internal) | Sets up [`charmbracelet/log`](https://github.com/charmbracelet/log) to make logging less boring                                                                     |
| `rhttp` (internal)  | Wrapper over `net/http` with provided backoff and retry policies that can be customised                                                                             |

## Usage

Examples of how to use the crawler package can be found in the [`example`](https://github.com/yusufaine/gocrawler/tree/main/example) directory.
