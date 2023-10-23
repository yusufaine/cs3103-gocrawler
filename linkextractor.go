package gocrawler

import (
	"bytes"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

var URLRegex = regexp.MustCompile(`[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,24}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)

// Takes in a map of blacklisted hosts and the response body and returns a slice of links
type LinkExtractor func(bl map[string]struct{}, currLink string, resp []byte) []string

// DefaultLinkExtractor looks for <a href="..."> tags and extracts the link if the host
// is not blacklisted.
func DefaultLinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		log.Error("unable to parse response body", "error", err)
		return nil
	}

	linkSet := make(map[string]struct{})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if !ok {
			return
		}

		// skip if link cannot be parsed
		newUrl, err := url.Parse(strings.TrimSpace(link))
		if err != nil {
			return
		}

		if newUrl.Host == "" || newUrl.Scheme == "" {
			return
		}

		// skip if host is blacklisted
		if _, ok := bl[newUrl.Host]; ok {
			return
		}

		// skip if link is not a valid URL
		if !URLRegex.MatchString(newUrl.String()) {
			return
		}

		linkSet[newUrl.String()] = struct{}{}
	})

	links := make([]string, 0, len(linkSet))
	for k := range linkSet {
		links = append(links, k)
	}
	slices.Sort(links)

	return links
}
