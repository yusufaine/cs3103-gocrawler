package gocrawler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
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
	PageMutex       sync.RWMutex
}

// New creates a new crawler client using the context to allow for cancellation, the crawler
// config, and list of response matchers to filter out responses.
//
// Note that the ordering of the response matchers matter, the first matcher to return
// false will cause the link to be skipped.
func New(ctx context.Context, config *Config, rm []ResponseMatcher) *Client {
	if len(rm) == 0 {
		rm = []ResponseMatcher{IsNoopResponse}
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
		rm:              rm,
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
func (c *Client) Crawl(ctx context.Context, le LinkExtractor, currDepth int, currLink, parent string) {

	// skip if the URL has been visited
	c.PageMutex.Lock()
	_, ok := c.VisitedPageInfo[currLink]
	c.PageMutex.Unlock()
	if ok {
		return
	}

	log.Info("visiting", "depth", currDepth, "link", currLink)

	resp := c.extractResponseBody(currLink, currDepth)
	if resp == nil {
		return
	}

	// returns a list of links whose hosts are not in the blacklist
	links := le(c.HostBlacklist, currLink, resp)

	// mark the current URL as visited
	c.PageMutex.Lock()
	c.VisitedPageInfo[currLink] = PageInfo{
		Content: resp,
		Depth:   currDepth,
		Links:   links,
		Parent:  parent,
	}
	c.PageMutex.Unlock()

	// crawl all outgoing links concurrently
	nextDepth := currDepth + 1
	var wg sync.WaitGroup
	for _, nextLink := range links {
		wg.Add(1)
		go func(currLink, nextLink string, nextDepth int) {
			defer wg.Done()
			// Do not continue crawling if the nextDepth has exceeded the max depth
			if nextDepth > c.MaxDepth {
				return
			}

			c.PageMutex.RLock()
			_, ok := c.VisitedPageInfo[nextLink]
			c.PageMutex.RUnlock()
			if ok {
				return
			}

			// ensure RPS is enforced
			_ = c.rl.Wait(ctx)
			c.Crawl(ctx, le, nextDepth, nextLink, currLink)
		}(currLink, nextLink, nextDepth)
	}
	wg.Wait()
	log.Info("visited all links in the branch", "depth", currDepth, "branch", currLink)
}

// Does the actual HTTP GET request and returns the response body if the response is
// successful and the content type is text.
func (c *Client) extractResponseBody(link string, depth int) []byte {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		log.Error("unable to parse url", "url", link, "error", err)
		return nil
	}

	remoteAddrs, err := net.LookupIP(parsedUrl.Host)
	if err != nil {
		log.Error("unable to resolve host", "host", parsedUrl.Host, "error", err)
		return nil
	}

	reqStart := time.Now()
	req, err := http.NewRequestWithContext(c.ctx, "GET", parsedUrl.String(), nil)
	if err != nil {
		log.Error("unable to create request", "url", parsedUrl.String(), "error", err)
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
		log.Error("unable to get response", "host", parsedUrl.Host, "error", err)
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
		remoteIpInfo := make([]IPInfo, 0, len(remoteAddrs))
		for _, addr := range remoteAddrs {
			asn, location, err := c.resolveIPInfo(addr.String())
			if err != nil {
				log.Warn("unable to resolve ip location", "error", err)
			}
			remoteIpInfo = append(remoteIpInfo, IPInfo{
				IP:       addr.String(),
				Location: location,
				ASNumber: asn,
			})
		}
		c.VisitedNetInfo[parsedUrl.Host] = []NetworkInfo{
			{
				RemoteIPInfo:        remoteIpInfo,
				VisitedPathSet:      map[string]struct{}{parsedUrl.Path: {}},
				TotalResponseTimeMs: respTime.Milliseconds(),
			},
		}
	}
	c.NetMutex.Unlock()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("unable to read response body", "url", link, "error", err)
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
