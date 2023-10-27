package gocrawler

// These values can be used for the users' benefit should they want to pass it
// to another program or export it to a JSON file for convenience.

type IPInfo struct {
	IP       string `json:"ip"`
	Location string `json:"location"`
	ASNumber string `json:"as_number"`
}

type NetworkInfo struct {
	RemoteIPInfo  []IPInfo `json:"remote_ip_info"`
	AvgResponseMs int64    `json:"avg_response_ms"`
	PathCount     int      `json:"path_count"`
	VisitedPaths  []string `json:"visited_paths"`

	// These values are not exported to JSON
	TotalResponseTimeMs int64               `json:"-"`
	VisitedPathSet      map[string]struct{} `json:"-"`
}

type PageInfo struct {
	Depth  int      `json:"depth"`
	Parent string   `json:"parent"`
	Links  []string `json:"links"`

	// These values are not exported to JSON
	Content []byte `json:"-"`
}
