package keywordcounter

import (
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/rhttp"
)

type toVisitInfo struct {
	depth int
	link  string
}

var (
	metadata = make(Metadata)
	metaMut  = sync.Mutex{}
	visited  = make(map[string]struct{})
	visitMut = sync.Mutex{}
	toVisit  = make([]toVisitInfo, 0)
)

// TODO: make concurrent
func Run(config *Config) {
	// TODO: Figure out best way to do BFS/DFS

	// Dummy code
	for _, v := range toVisit {
		if v.depth > config.MaxDepth {
			continue
		}

		visitMut.Lock()
		_, ok := visited[v.link]
		if ok {
			visitMut.Unlock()
			continue
		}
		visitMut.Unlock()

		toEnqueue := visit(config, v.link, v.depth)
		// only mark as visited after resolving
		visitMut.Lock()
		visited[v.link] = struct{}{}
		visitMut.Unlock()
		toVisit = append(toVisit, toEnqueue...)
	}

	metadata.ExportAsJSON(config.OutputFile)
}

func visit(config *Config, link string, depth int) []toVisitInfo {
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Error("failed to parse url", "error", err.Error())
	}

	if _, ok := config.BlacklistHosts[parsedURL.Host]; ok {
		log.Warn("skipping blacklisted host", "host", parsedURL.Host)
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		log.Error("failed to create request", "error", err.Error())
		return nil
	}

	resp, err := rhttp.New().Do(req)
	if err != nil {
		log.Error("failed to get response", "error", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("received non-200 response", "status", resp.StatusCode, "url", link)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", "error", err.Error())
		return nil
	}

	keywords := extractKeywords(string(body), []string{})
	links := extractLinks(string(body))

	// TODO: shift this to a separate function
	metaMut.Lock()
	defer metaMut.Unlock()
	entry, ok := metadata[link]
	if !ok {
		entry = struct {
			NetworkInfo
			Pages []PageInfo
		}{
			NetworkInfo: NetworkInfo{
				RemoteAddr: req.RemoteAddr,
				Path:       parsedURL.Path,
			},
			Pages: []PageInfo{
				{
					Depth:       depth,
					Links:       links,
					KeywordFreq: keywords,
				},
			},
		}
	} else {
		entry.Pages = append(entry.Pages, PageInfo{
			Depth:       depth,
			Links:       links,
			KeywordFreq: keywords,
		})
	}

	// Run's responsibility to ensure correct visiting
	toEnqueue := make([]toVisitInfo, len(links))
	for _, l := range links {
		toEnqueue = append(toEnqueue, toVisitInfo{
			depth: depth + 1,
			link:  l,
		})
	}
	return toEnqueue
}

// TODO: seems inefficient to count keywords N-times
func extractKeywords(body string, keywords []string) map[string]int {
	// TODO: body preprocessing if needed
	kwMap := make(map[string]int)
	for _, kw := range keywords {
		kwMap[kw] = strings.Count(body, kw)
	}
	panic("extractKeywords not implemented")

	return kwMap
}

// TODO: extract links from body
func extractLinks(body string) []string {
	// TODO: body preprocessing if needed
	linkSet := make(map[string]struct{})

	// TODO: extract links from body, regex or some other method
	var links []string
	for _, l := range links {
		linkSet[l] = struct{}{}
	}
	links = make([]string, 0, len(linkSet))
	for k := range linkSet {
		links = append(links, k)
	}
	slices.Sort(links)
	panic("extractLinks not implemented")

	return links
}
