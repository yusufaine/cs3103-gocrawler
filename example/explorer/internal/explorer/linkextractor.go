package explorer

import (
	"net/url"
	"strings"
	"sync"

	"github.com/yusufaine/gocrawler"
)

// Extracts links that are not blacklisted, and has the scheme "http" or "https"
func ExplorerLinkExtractor(c *gocrawler.Client, currLink string, resp []byte) []string {
	links := gocrawler.DefaultLinkExtractor(c, currLink, resp)
	blHosts := make([]string, 0, len(c.HostBlacklist))
	for k := range c.HostBlacklist {
		blHosts = append(blHosts, k)
	}

	var (
		wg            sync.WaitGroup
		filtered      []string
		filteredMutex sync.Mutex
	)
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			defer wg.Done()

			p, err := url.Parse(link)
			if err != nil {
				return
			}

			for _, host := range blHosts {
				if strings.Contains(p.Host, host) || strings.Contains(host, p.Host) {
					return
				}
			}

			if p.Scheme != "https" && p.Scheme != "http" {
				return
			}

			filteredMutex.Lock()
			defer filteredMutex.Unlock()
			filtered = append(filtered, link)
		}(link)
	}
	wg.Wait()

	return filtered
}
