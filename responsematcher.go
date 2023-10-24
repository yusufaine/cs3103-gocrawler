package gocrawler

import (
	"net/http"
	"strings"
)

// ResponseMatcher is a function that takes an http.Response and returns a boolean to indicate
// whether or not the contents of the URL should be processed (e.g extract links)
type ResponseMatcher func(resp *http.Response) bool

// This matches all responses
func IsNoopResponse(resp *http.Response) bool {
	return true
}

// This matches all responses that return a 200 status code
func IsOkResponse(resp *http.Response) bool {
	return resp.StatusCode == 200
}

// This matches all responses that return a 4xx status code
func IsClientErrorResponse(resp *http.Response) bool {
	return resp.StatusCode >= 400 && resp.StatusCode < 500
}

// This matches all responses that return a 5xx status code
func IsServerErrorResponse(resp *http.Response) bool {
	return resp.StatusCode >= 500
}

// This matches all responses that return a 2xx status code and have a Content-Type header
// that contains "text/html"
func IsHtmlContent(resp *http.Response) bool {
	return IsOkResponse(resp) && strings.Contains(resp.Header.Get("Content-Type"), "text/html")
}
