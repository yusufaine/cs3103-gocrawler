package tianalyser

import (
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/yusufaine/gocrawler"
)

// Returns a list of all outgoing links with the same host as the current link and the path contains
// "dota2/The_International/"
func TILinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
	var (
		filteredLinks []string
		filterMutex   sync.Mutex
		wg            sync.WaitGroup
	)

	currURL, err := url.Parse(currLink)
	if err != nil {
		return nil
	}

	allLinks := gocrawler.DefaultLinkExtractor(bl, currLink, resp)
	wg.Add(len(allLinks))
	for _, link := range allLinks {
		go func(link string) {
			defer wg.Done()
			toFilterURL, err := url.Parse(link)
			if err != nil {
				return
			}

			// skip if host is not liquipedia.net
			if toFilterURL.Host != currURL.Host {
				return
			}

			// skip if path does not contain "dota2/The_International/"
			if !strings.Contains(toFilterURL.Path, "dota2/The_International/") {
				return
			}

			filterMutex.Lock()
			defer filterMutex.Unlock()
			filteredLinks = append(filteredLinks, link)
		}(link)
	}
	wg.Wait()

	slices.Sort(filteredLinks)

	return filteredLinks
}
