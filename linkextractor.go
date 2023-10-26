package gocrawler

import (
	"bytes"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

// Takes in a map of blacklisted hosts and the response body and returns a slice of links
type LinkExtractor func(c *Client, currLink string, resp []byte) []string

// DefaultLinkExtractor looks for <a href="..."> tags and extracts the link if the host
// is not blacklisted. This function assumes that if the href value is a relative path,
// it is relative to the current URL.
func DefaultLinkExtractor(c *Client, currLink string, resp []byte) []string {
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

		// if host is empty, likely to be a relative path
		if outURL.Host == "" && outURL.Path != "" {
			tmpPath := outURL.Path
			outURL = currURL
			outURL.Path = tmpPath
		}

		// skip if host is blacklisted
		if _, ok := c.HostBlacklist[outURL.Host]; ok {
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
