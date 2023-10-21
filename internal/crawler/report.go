package crawler

import (
	"slices"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/cs3203-g46-crawler/pkg/filewriter"
)

type NetworkInfo struct {
	VisitedPaths        []string `json:"paths"`
	RemoteAddrs         []string `json:"remote_addr"`
	DNSAddrs            []string `json:"dns_addrs"`
	TotalResponseTimeMs int64    `json:"-"`
	AvgResponseMs       int64    `json:"avg_response_ms"`
}

type PageInfo struct {
	Content []byte   `json:"-"`
	Depth   int      `json:"depth"`
	Links   []string `json:"links"`
}

type ReportFormat struct {
	Seed      string   `json:"seed"`
	Depth     int      `json:"max_depth"`
	Blacklist []string `json:"blacklist"`

	VisitedNetInfo  map[string][]NetworkInfo `json:"network_info"`
	VisitedPageResp map[string]PageInfo      `json:"page_info"`
}

func (cr *Crawler) GenerateReport(config *Config) {
	bls := make([]string, 0, len(cr.HostBlacklist))
	for k := range cr.HostBlacklist {
		bls = append(bls, k)
	}
	slices.Sort(bls)

	report := ReportFormat{
		Seed:            config.SeedURL.String(),
		Depth:           config.MaxDepth,
		Blacklist:       bls,
		VisitedNetInfo:  cr.VisitedNetInfo,
		VisitedPageResp: cr.VisitedPageInfo,
	}
	for k, v := range report.VisitedNetInfo {
		for i, v1 := range v {
			slices.Sort(v1.DNSAddrs)
			slices.Sort(v1.RemoteAddrs)
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

	if err := filewriter.ToJSON(report, config.RelReportPath); err != nil {
		log.Error("unable to write to file",
			"file", config.RelReportPath,
			"error", err)
	} else {
		log.Info("exported crawler report", "file", config.RelReportPath)
	}
}
