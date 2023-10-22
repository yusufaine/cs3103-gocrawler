package gocrawler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler/internal/rhttp"
	"golang.org/x/time/rate"
)

type Client struct {
	ctx context.Context
	hc  *rhttp.Client
	rl  *rate.Limiter
	rm  []ResponseMatcher

	MaxDepth        int
	HostBlacklist   map[string]struct{}
	VisitedNetInfo  map[string][]NetworkInfo
	NetMutex        sync.Mutex
	VisitedPageInfo map[string]PageInfo
	PageMutex       sync.Mutex
}

// Creates a new crawler client using the context to allw for cancellation, the crawler
// config, and list of response matchers to filter out responses.
//
// Note that the ordering of the response matchers matter, the first matcher to return
// false will cause the link to be skipped.
func New(ctx context.Context, config *Config, rm []ResponseMatcher) *Client {
	if len(rm) == 0 {
		rm = []ResponseMatcher{NoopResponseFilter}
		log.Warn("no response matchers supplied, accepting all responses")
	}

	retryClient := rhttp.New(
		rhttp.WithBackoffPolicy(rhttp.ExponentialBackoff),
		rhttp.WithMaxRetries(config.MaxRetries),
		rhttp.WithRetryPolicy(rhttp.DefaultRetry),
		rhttp.WithTimeout(config.Timeout),
		rhttp.WithProxy(config.ProxyURL),
	)

	c := &Client{
		ctx:             ctx,
		hc:              retryClient,
		rl:              rate.NewLimiter(rate.Limit(config.MaxRPS), 1),
		MaxDepth:        config.MaxDepth - 1,
		HostBlacklist:   config.BlacklistHosts,
		VisitedNetInfo:  make(map[string][]NetworkInfo),
		VisitedPageInfo: make(map[string]PageInfo),
	}

	return c
}

// Crawl is called recursively to crawl the supplied URL and all outgoing links which is
// extracted by the supplied LinkExtractor. The crawl will stop when the MaxDepth is reached
// or if the context is cancelled.
func (c *Client) Crawl(ctx context.Context, le LinkExtractor, parsedURL *url.URL, currDepth int) {
	// skip if the current depth is greater than the max depth
	if currDepth > c.MaxDepth {
		return
	}

	// skip if the URL has been visited
	c.PageMutex.Lock()
	if _, ok := c.VisitedPageInfo[parsedURL.String()]; ok {
		c.PageMutex.Unlock()
		return
	}
	c.PageMutex.Unlock()

	log.Info("visiting", "depth", currDepth, "link", parsedURL.String())

	resp := c.extractResponseBody(parsedURL.String(), currDepth)
	if resp == nil {
		return
	}

	// returns a list of links whose hosts are not in the blacklist
	links := le(c, parsedURL, resp)
	strLinks := make([]string, 0, len(links))
	for _, l := range links {
		strLinks = append(strLinks, l.String())
	}

	// mark the current URL as visited
	c.PageMutex.Lock()
	c.VisitedPageInfo[parsedURL.String()] = PageInfo{
		Content: resp,
		Depth:   currDepth,
		Links:   strLinks,
	}
	c.PageMutex.Unlock()

	// crawl all outgoing links concurrently
	currDepth++
	var wg sync.WaitGroup
	for _, l := range links {
		wg.Add(1)
		go func(l *url.URL, currDepth int) {
			defer wg.Done()
			// ensure RPS is enforced
			_ = c.rl.Wait(ctx)
			c.Crawl(ctx, le, l, currDepth)
		}(l, currDepth)
	}
	wg.Wait()
}

// Does the actual HTTP GET request and returns the response body if the response is
// successful and the content type is text.
func (c *Client) extractResponseBody(link string, depth int) []byte {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		log.Error("unable to parse url",
			"url", link,
			"error", err)
		return nil
	}

	var remoteAddr string
	reqStart := time.Now()
	req, err := http.NewRequestWithContext(c.ctx, "GET", parsedUrl.String(), nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr().String()
		},
	}))
	if err != nil {
		log.Error("unable to create request",
			"url", parsedUrl.String(),
			"error", err)
		return nil
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		c.PageMutex.Lock()
		c.VisitedPageInfo[link] = PageInfo{Depth: depth}
		defer c.PageMutex.Unlock()
		log.Error("unable to get response",
			"host", parsedUrl.Host,
			"error", err)
		return nil
	}
	defer resp.Body.Close()

	respTime := time.Since(reqStart)

	// if any of the response filters return false, skip the link
	for _, f := range c.rm {
		if !f(resp) {
			return nil
		}
	}

	c.NetMutex.Lock()
	if infos, ok := c.VisitedNetInfo[parsedUrl.Host]; ok {
		for i, info := range infos {
			if _, ok := info.VisitedPathSet[parsedUrl.Path]; !ok {
				info.VisitedPathSet[parsedUrl.Path] = struct{}{}
			}

			info.TotalResponseTimeMs += respTime.Milliseconds()
			c.VisitedNetInfo[parsedUrl.Host][i] = info
		}
	} else {
		asn, location, err := c.resolveIPInfo(strings.Split(remoteAddr, ":")[0])
		if err != nil {
			log.Warn("unable to resolve ip location", "error", err)
		}
		c.VisitedNetInfo[parsedUrl.Host] = []NetworkInfo{
			{
				RemoteAddr:          remoteAddr,
				Location:            location,
				ASNumber:            asn,
				VisitedPathSet:      map[string]struct{}{parsedUrl.Path: {}},
				TotalResponseTimeMs: respTime.Milliseconds(),
			},
		}
	}
	c.NetMutex.Unlock()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("unable to read response body",
			"url", link,
			"error", err)
		return nil
	}
	return body
}

func (c *Client) resolveIPInfo(ip string) (string, string, error) {
	req, err := http.NewRequestWithContext(c.ctx, "GET", "https://ipapi.co/"+ip+"/json/", nil)
	if err != nil {
		return "", "", err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var ipInfo struct {
		ASNumber string `json:"asn,omitempty"`
		Country  string `json:"country_name,omitempty"`
		Region   string `json:"region,omitempty"`
	}
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", "", err
	}

	return ipInfo.ASNumber, ipInfo.Country + ", " + ipInfo.Region, nil
}
