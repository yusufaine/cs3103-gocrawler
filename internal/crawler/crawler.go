package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
	"golang.org/x/time/rate"
)

type Crawler struct {
	ctx context.Context
	hc  *rhttp.Client
	rl  *rate.Limiter

	MaxDepth        int
	HostBlacklist   map[string]struct{}
	VisitedNetInfo  map[string][]NetworkInfo
	NetMutex        sync.Mutex
	VisitedPageInfo map[string]PageInfo
	PageMutex       sync.Mutex
}

// To blacklist remote hosts, use WithBlacklist()
func New(ctx context.Context, config *Config, maxRPS float64) *Crawler {

	retryClient := rhttp.New(
		rhttp.WithBackoffPolicy(rhttp.DefaultLinearBackoff),
		rhttp.WithMaxRetries(config.MaxRetries),
		rhttp.WithRetryPolicy(rhttp.DefaultRetry),
		rhttp.WithTimeout(config.Timeout),
	)

	c := &Crawler{
		ctx:             ctx,
		hc:              retryClient,
		rl:              rate.NewLimiter(rate.Limit(maxRPS), 1),
		MaxDepth:        config.MaxDepth - 1,
		HostBlacklist:   config.BlacklistHosts,
		VisitedNetInfo:  make(map[string][]NetworkInfo),
		VisitedPageInfo: make(map[string]PageInfo),
	}

	return c
}

func (c *Crawler) Crawl(ctx context.Context, le LinkExtractor, parsedURL *url.URL, currDepth int) {
	c.PageMutex.Lock()
	if _, ok := c.VisitedPageInfo[parsedURL.String()]; ok {
		c.PageMutex.Unlock()
		return
	}

	if currDepth > c.MaxDepth {
		c.PageMutex.Unlock()
		return
	}
	c.PageMutex.Unlock()

	log.Info("visiting", "depth", currDepth, "link", parsedURL.String())

	resp := c.extractResponseBody(parsedURL.String(), currDepth)
	if resp == nil {
		return
	}

	links := le(c, parsedURL, resp)
	strLinks := make([]string, 0, len(links))
	for _, l := range links {
		strLinks = append(strLinks, l.String())
	}
	c.PageMutex.Lock()
	c.VisitedPageInfo[parsedURL.String()] = PageInfo{
		Content: resp,
		Depth:   currDepth,
		Links:   strLinks,
	}
	c.PageMutex.Unlock()

	currDepth++
	var wg sync.WaitGroup
	for _, l := range links {
		wg.Add(1)
		go func(l *url.URL, currDepth int) {
			defer wg.Done()
			_ = c.rl.Wait(ctx) // ignore error
			c.Crawl(ctx, le, l, currDepth)
		}(l, currDepth)
	}
	wg.Wait()
}

func (c *Crawler) extractResponseBody(link string, depth int) []byte {
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
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "text") {
		return nil
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

func (c *Crawler) resolveIPInfo(ip string) (string, string, error) {
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

	return ipInfo.ASNumber, fmt.Sprintf("%s, %s", ipInfo.Country, ipInfo.Region), nil
}
