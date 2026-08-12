package model

// Struct made to fix routes with only one function
type LineInfo struct {
	Line          string `json:"line"`
	Destination   string `json:"destination"`
	Journey_code  string `json:"journey_code"`
	Official_line string `json:"official_line"`
	Line_type     string `json:"line_type"`
}
