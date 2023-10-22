package crawler

//* These values can be used for the users' benefit should they want to pass it
//* to another program or export it to a JSON file for convenience.

type NetworkInfo struct {
	VisitedPaths  []string `json:"paths"`
	RemoteAddr    string   `json:"remote_addr"`
	Location      string   `json:"location"`
	ASNumber      string   `json:"as_number"`
	AvgResponseMs int64    `json:"avg_response_ms"`
	PathCount     int      `json:"path_count"`

	// These values are not exported to JSON
	TotalResponseTimeMs int64               `json:"-"`
	VisitedPathSet      map[string]struct{} `json:"-"`
}

type PageInfo struct {
	Depth  int      `json:"depth"`
	Links  []string `json:"links"`
	Parent string   `json:"parent"`

	// These values are not exported to JSON
	Content []byte `json:"-"`
}
