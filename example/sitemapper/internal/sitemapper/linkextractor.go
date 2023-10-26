package sitemapper

import (
	"net/url"
	"slices"
	"sync"

	"github.com/yusufaine/gocrawler"
)

func SameHostLinkExtractor(c *gocrawler.Client, currLink string, resp []byte) []string {
	var (
		filteredLinks []string
		filterMutex   sync.Mutex
		wg            sync.WaitGroup
	)
	currURL, err := url.Parse(currLink)
	if err != nil {
		return nil
	}

	allLinks := gocrawler.DefaultLinkExtractor(c, currLink, resp)
	wg.Add(len(allLinks))
	for _, link := range allLinks {
		go func(link string) {
			defer wg.Done()
			parsedURL, err := url.Parse(link)
			if err != nil {
				return
			}

			if parsedURL.Host != currURL.Host {
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
