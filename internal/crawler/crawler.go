package crawler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
	"golang.org/x/time/rate"
)

type Crawler struct {
	ctx      context.Context
	hc       *rhttp.Client
	rl       *rate.Limiter
	dnsMutex sync.Mutex
	netMutex sync.Mutex

	MaxDepth        int
	HostBlacklist   map[string]struct{}
	VisitedNetInfo  map[string][]NetworkInfo
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

	var (
		remoteAddr string
		dnsAddrs   []string
	)
	reqStart := time.Now()
	req, err := http.NewRequestWithContext(c.ctx, "GET", parsedUrl.String(), nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr().String()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			c.dnsMutex.Lock()
			defer c.dnsMutex.Unlock()
			for _, addr := range dnsInfo.Addrs {
				dnsAddrs = append(dnsAddrs, addr.String())
			}
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

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "text") {
		return nil
	}
	respTime := time.Since(reqStart)

	// deduplicate dns addresses
	c.dnsMutex.Lock()
	dnsSet := make(map[string]struct{})
	for _, da := range dnsAddrs {
		dnsSet[da] = struct{}{}
	}
	dnsAddrs = make([]string, 0, len(dnsSet))
	for k := range dnsSet {
		dnsAddrs = append(dnsAddrs, k)
	}
	c.dnsMutex.Unlock()

	c.netMutex.Lock()
	if infos, ok := c.VisitedNetInfo[parsedUrl.Host]; ok {
		for i, info := range infos {
			// host may have multiple remote addresses and DNS addresses
			if !slices.Contains(info.RemoteAddrs, remoteAddr) {
				info.RemoteAddrs = append(info.RemoteAddrs, remoteAddr)
			}

			for _, da := range dnsAddrs {
				if !slices.Contains(info.DNSAddrs, da) {
					info.DNSAddrs = append(info.DNSAddrs, da)
				}
			}

			if !slices.Contains(info.VisitedPaths, parsedUrl.Path) {
				info.VisitedPaths = append(info.VisitedPaths, parsedUrl.Path)
			}

			info.TotalResponseTimeMs += respTime.Milliseconds()
			c.VisitedNetInfo[parsedUrl.Host][i] = info
		}
	} else {
		c.VisitedNetInfo[parsedUrl.Host] = []NetworkInfo{
			{
				RemoteAddrs:         []string{remoteAddr},
				VisitedPaths:        []string{parsedUrl.Path},
				DNSAddrs:            dnsAddrs,
				TotalResponseTimeMs: respTime.Milliseconds(),
			},
		}
	}
	c.netMutex.Unlock()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("unable to read response body",
			"url", link,
			"error", err)
		return nil
	}
	return body
}
