package sitemapper

import (
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler"
	"github.com/yusufaine/gocrawler/example/internal/filewriter"
)

type ReportFormat struct {
	Seed      string  `json:"seed"`
	MaxRPS    float64 `json:"max_rps"`
	CrawlTime string  `json:"crawl_time"`

	VisitedNetInfo  map[string][]gocrawler.NetworkInfo `json:"network_info"`
	VisitedPageResp map[string]gocrawler.PageInfo      `json:"page_info"`
}

// Generates a report in JSON format from the crawler client and config. The report contains
// the initial crawler info, the network info for each host visited, and the page info for
// each page visited such as all the links found in the page if the link belongs to the same
// host as the seed URL.
func Generate(config *Config, cr *gocrawler.Client, elapsed time.Duration) {
	report := ReportFormat{
		Seed:            config.SeedURLs[0],
		MaxRPS:          config.MaxRPS,
		CrawlTime:       elapsed.String(),
		VisitedNetInfo:  cr.VisitedNetInfo,
		VisitedPageResp: cr.VisitedPageInfo,
	}
	for k, v := range report.VisitedNetInfo {
		for i, v1 := range v {
			v1.PathCount = len(v1.VisitedPathSet)
			v1.VisitedPaths = make([]string, 0, v1.PathCount)
			for k := range v1.VisitedPathSet {
				v1.VisitedPaths = append(v1.VisitedPaths, k)
			}
			slices.Sort(v1.VisitedPaths)
			report.VisitedNetInfo[k][i] = v1

			visitedCount := int64(len(v1.VisitedPaths))
			if visitedCount == 0 {
				visitedCount = 1
			}
			v1.AvgResponseMs = v1.TotalResponseTimeMs / visitedCount
			report.VisitedNetInfo[k][i] = v1
		}
	}

	if err := filewriter.ToJSON(report, config.ReportPath); err != nil {
		log.Error("unable to write to file", "file", config.ReportPath, "error", err)
	} else {
		log.Info("exported crawler report", "file", config.ReportPath)
	}
}
