package tianalyser

import (
	"bytes"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

// Returns a list of all outgoing links from the page
func ReportLinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
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

		if !strings.Contains(newUrl.Path, "dota2/The_International/") {
			return
		}

		updatedURL, err := url.Parse(currLink)
		if err != nil {
			return
		}
		updatedURL.Path = newUrl.Path

		linkSet[updatedURL.String()] = struct{}{}
	})

	links := make([]string, 0, len(linkSet))
	for k := range linkSet {
		links = append(links, k)
	}
	slices.Sort(links)

	return links
}

// Returns a list of all outgoing links that have not been visited from any page
func TILinkExtractor(bl map[string]struct{}, currLink string, resp []byte) []string {
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

		if !strings.Contains(newUrl.Path, "dota2/The_International/") {
			return
		}

		updatedURL, err := url.Parse(currLink)
		if err != nil {
			return
		}
		updatedURL.Path = newUrl.Path

		linkSet[updatedURL.String()] = struct{}{}
	})

	links := make([]string, 0, len(linkSet))
	for k := range linkSet {
		links = append(links, k)
	}
	slices.Sort(links)

	return links
}
