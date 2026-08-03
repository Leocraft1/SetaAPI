package arrivals

type ArrivalRaw struct {
	Arrival struct {
        Services []ServiceRaw `json:"services"`
    } `json:"arrival"`
}