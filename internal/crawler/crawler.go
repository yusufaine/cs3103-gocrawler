package crawler

import (
	"context"
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

type NetworkInfo struct {
	VisitedPaths   []string `json:"paths"`
	RemoteAddrs    []string `json:"remote_addr"`
	DNSAddrs       []string `json:"dns_addrs"`
	ResponseTimeMs int64    `json:"response_ms"`
}

type PageInfo struct {
	Content []byte   `json:"-"`
	Depth   int      `json:"depth"`
	Links   []string `json:"links"`
}

type Crawler struct {
	le              LinkExtractor
	hc              *rhttp.Client
	rl              *rate.Limiter
	MaxDepth        int
	HostBlacklist   map[string]struct{}
	VisitedNetInfo  map[string][]NetworkInfo
	VisitedPageResp map[string]PageInfo
}

// To blacklist remote hosts, use WithBlacklist()
func New(maxDepth int, opts ...CrawlerOption) *Crawler {
	c := &Crawler{
		le:              DefaultLinkExtractor,
		hc:              rhttp.New(rhttp.WithTimeout(3 * time.Second)),
		MaxDepth:        maxDepth - 1,
		HostBlacklist:   make(map[string]struct{}),
		VisitedNetInfo:  make(map[string][]NetworkInfo),
		VisitedPageResp: make(map[string]PageInfo),
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Crawler) Crawl(ctx context.Context, link string, currDepth int) {
	if _, ok := c.VisitedPageResp[link]; ok {
		return
	}

	if currDepth > c.MaxDepth {
		return
	}

	log.Info("visiting", "depth", currDepth, "link", link)

	resp := c.extractResponseBody(link, currDepth)
	if resp == nil {
		return
	}

	links := c.le(c.HostBlacklist, resp)
	c.VisitedPageResp[link] = PageInfo{
		Content: resp,
		Depth:   currDepth,
		Links:   links,
	}

	currDepth++
	var wg sync.WaitGroup
	for _, l := range links {
		wg.Add(1)
		go func(l string, d int) {
			defer wg.Done()
			_ = c.rl.Wait(ctx) // ignore error
			c.Crawl(ctx, l, currDepth)
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
	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr().String()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			for _, addr := range dnsInfo.Addrs {
				if !slices.Contains(dnsAddrs, addr.String()) {
					dnsAddrs = append(dnsAddrs, addr.String())
				}
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
		log.Error("unable to get response",
			"url", parsedUrl.String(),
			"error", err)
		return nil
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text") || contentType == "" {
		log.Debug("skipping non-text response",
			"type", resp.Header.Get("Content-Type"),
			"link", link)
		return nil
	}
	respTime := time.Since(reqStart)

	if infos, ok := c.VisitedNetInfo[parsedUrl.Host]; ok {
		for _, info := range infos {
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

			info.ResponseTimeMs = respTime.Milliseconds()
		}
	} else {
		c.VisitedNetInfo[parsedUrl.Host] = []NetworkInfo{
			{
				RemoteAddrs:    []string{remoteAddr},
				VisitedPaths:   []string{parsedUrl.Path},
				DNSAddrs:       dnsAddrs,
				ResponseTimeMs: respTime.Milliseconds(),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("unable to read response body",
			"url", link,
			"error", err)
		return nil
	}
	return body
}
