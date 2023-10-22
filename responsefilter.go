package gocrawler

import (
	"net/http"
	"strings"
)

type ResponseMatcher func(resp *http.Response) bool

// This matches all responses
func NoopResponseFilter(resp *http.Response) bool {
	return true
}

// This matches all responses that return a 200 status code
func OkResponseFilter(resp *http.Response) bool {
	return resp.StatusCode == 200
}

// This matches all responses that return a 4xx status code
func ClientErrorResponseFilter(resp *http.Response) bool {
	return resp.StatusCode >= 400 && resp.StatusCode < 500
}

// This matches all responses that return a 5xx status code
func ServerErrorResponseFilter(resp *http.Response) bool {
	return resp.StatusCode >= 500
}

// This matches all responses that return a 2xx status code and have a Content-Type header
// that contains "text/html"
func HtmlContentFilter(resp *http.Response) bool {
	return OkResponseFilter(resp) && strings.Contains(resp.Header.Get("Content-Type"), "text/html")
}
