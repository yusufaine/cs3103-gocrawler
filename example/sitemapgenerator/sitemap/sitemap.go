package sitemap

import (
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/gocrawler"
	"github.com/yusufaine/gocrawler/example/internal/filewriter"
)

type ReportFormat struct {
	Seeds     []string `json:"seeds"`
	Depth     int      `json:"max_depth"`
	Blacklist []string `json:"blacklist"`
	MaxRPS    float64  `json:"max_rps"`
	CrawlTime string   `json:"crawl_time"`

	VisitedNetInfo  map[string][]gocrawler.NetworkInfo `json:"network_info"`
	VisitedPageResp map[string]gocrawler.PageInfo      `json:"page_info"`
}

func Generate(config *Config, cr *gocrawler.Client, elapsed time.Duration) {
	bls := make([]string, 0, len(cr.HostBlacklist))
	for k := range cr.HostBlacklist {
		bls = append(bls, k)
	}
	slices.Sort(bls)

	report := ReportFormat{
		Seeds:           config.SeedURLs,
		Depth:           config.MaxDepth,
		Blacklist:       bls,
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
