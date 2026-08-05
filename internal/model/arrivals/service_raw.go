package arrivals

type ServiceRaw struct {
	Service         string  `json:"service"`
	Arrival         string  `json:"arrival"`
	Type            string  `json:"type"`
	Destination     string  `json:"destination"`
	FleetCode       string  `json:"fleetCode"`
	DutyId          string  `json:"dutyId"`
	Busnum          string  `json:"busnum"`
	ServiceType     string  `json:"serviceType"`
	OccupancyStatus *string `json:"occupancyStatus"`
	Codice_corsa    string  `json:"codice_corsa"`
	Posti_totali    *int    `json:"posti_totali"`
	Num_passeggeri  *int    `json:"num_passeggeri"`
	Next_stop       *string `json:"next_stop"`
}
