package arrivals

type Arrival struct {
	Arrival struct {
		Services []Service `json:"services"`
	} `json:"arrivals"`
}