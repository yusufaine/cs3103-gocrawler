package tianalyser

import (
	"bytes"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler"
	"github.com/yusufaine/gocrawler/example/internal/filewriter"
)

type CountryTableRow struct {
	Country        string   `json:"country"`
	Representation string   `json:"representation"`
	Players        []string `json:"players"`
}

type ReportFormat struct {
	Seed      string  `json:"seed"`
	MaxRPS    float64 `json:"max_rps"`
	CrawlTime string  `json:"crawl_time"`

	NetInfo map[string][]gocrawler.NetworkInfo `json:"network_info"`
	TIStats map[string][]CountryTableRow       `json:"ti_stats"`
}

// Generates a report in JSON format from the crawler client and config. The report contains
// the initial crawler info, the network info for each host visited, and the country representation
// table for each TI page visited.
func Generate(cr *gocrawler.Client, config *Config, elapsed time.Duration) {
	report := ReportFormat{
		Seed:      config.SeedURLs[0],
		MaxRPS:    config.MaxRPS,
		CrawlTime: elapsed.String(),
		NetInfo:   cr.VisitedNetInfo,
		TIStats:   make(map[string][]CountryTableRow),
	}
	for k, v := range report.NetInfo {
		for i, v1 := range v {
			v1.PathCount = len(v1.VisitedPathSet)
			v1.VisitedPaths = make([]string, 0, v1.PathCount)
			for k := range v1.VisitedPathSet {
				v1.VisitedPaths = append(v1.VisitedPaths, k)
			}
			slices.Sort(v1.VisitedPaths)
			report.NetInfo[k][i] = v1

			visitedCount := int64(len(v1.VisitedPaths))
			if visitedCount == 0 {
				visitedCount = 1
			}
			v1.AvgResponseMs = v1.TotalResponseTimeMs / visitedCount
			report.NetInfo[k][i] = v1
		}
	}

	links := make([]string, 0, len(cr.VisitedPageInfo))
	for k := range cr.VisitedPageInfo {
		links = append(links, k)
	}
	for _, l := range links {
		table := extractCountryRepresentationTable(cr.VisitedPageInfo[l].Content)
		if table != nil {
			report.TIStats[l] = table
		}
	}

	if err := filewriter.ToJSON(report, config.ReportPath); err != nil {
		log.Error("unable to write to file", "file", config.ReportPath, "error", err)
	} else {
		log.Info("exported DOTA TI report", "file", config.ReportPath)
	}
}

// Uses goquery to find the table containing the country representation in
// the response body, if any.
func extractCountryRepresentationTable(resp []byte) []CountryTableRow {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		panic(err)
	}

	var repTable []CountryTableRow
	doc.Find("#Country_Representation").
		Parent().Next().Each(func(i int, s *goquery.Selection) {
		var row CountryTableRow
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i % 4 {
			case 1:
				row.Country = strings.ReplaceAll(s.Text(), "\u00a0", "")
			case 2:
				row.Representation = s.Text()
			case 3:
				for _, pl := range strings.Split(s.Text(), ",") {
					pl = strings.TrimSpace(pl)
					if pl != "" {
						row.Players = append(row.Players, pl)
					}
				}
				repTable = append(repTable, row)
				row = CountryTableRow{}
			}
		})
	})

	return repTable
}
