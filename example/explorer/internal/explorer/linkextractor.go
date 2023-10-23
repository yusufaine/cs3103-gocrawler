package explorer

import (
	"bytes"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

func DepthLinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		log.Error("unable to parse response body", "error", err)
		return nil
	}

	currURL, err := url.Parse(currLink)
	if err != nil {
		return nil
	}

	linkSet := make(map[string]struct{})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		outLink, ok := s.Attr("href")
		if !ok {
			return
		}

		// skip if link cannot be parsed
		outURL, err := url.Parse(strings.TrimSpace(outLink))
		if err != nil {
			return
		}

		// if host is empty, likely to be a relative path (e.g /about)
		if outURL.Host == "" && outURL.Path != "" {
			outURL.Host = currURL.Host
		}

		// skip if host is blacklisted
		if _, ok := bl[outURL.Host]; ok {
			return
		}

		linkSet[outURL.String()] = struct{}{}
	})

	links := make([]string, 0, len(linkSet))
	for k := range linkSet {
		links = append(links, k)
	}
	slices.Sort(links)

	return links
}
