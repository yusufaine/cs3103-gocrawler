package keywordcounter

import (
	"encoding/json"

	"github.com/yusufaine/cs3203-g46-crawler/pkg/fileexporter"
)

type PageInfo struct {
	Depth       int            `json:"depth"` // depth from seed
	Links       []string       `json:"links"` // map of links found at current page
	KeywordFreq map[string]int `json:"keyword_freq"`
}

type NetworkInfo struct {
	RemoteAddr string `json:"remote_ip"` // remote IP address
	Path       string `json:"path"`      // path of the page
}

type Metadata map[string]struct {
	NetworkInfo
	Pages []PageInfo
}

func (m *Metadata) ExportAsJSON(filename string) error {
	d, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return fileexporter.WriteToFile(d, filename)
}
