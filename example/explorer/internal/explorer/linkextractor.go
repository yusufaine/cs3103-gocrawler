package explorer

import (
	"net/url"
	"strings"
	"sync"

	"github.com/yusufaine/gocrawler"
)

func ExplorerLinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
	links := gocrawler.DefaultLinkExtractor(bl, currLink, resp)
	blHosts := make([]string, 0, len(bl))
	for k := range bl {
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
