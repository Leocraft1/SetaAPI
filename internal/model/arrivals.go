package model

type Arrival struct {
	Arrival struct {
		Services []Service `json:"services"`
	} `json:"arrivals"`
}

type Service struct {
	LineInfo
	Arrival_time     string  `json:"arrival_time"`
	State            string  `json:"state"`
	Basin            string  `json:"basin"`
	Vehicle_table    string  `json:"vehicle_table"`
	Vehicle          string  `json:"vehicle"`
	Line_type        string  `json:"line_type"`
	Occupancy_status *string `json:"occupancy_status"`
	Total_room       *int    `json:"total_room"`
	Passenger_number *int    `json:"passenger_number"`
	Next_stop        *string `json:"next_stop"`
	Has_problems     bool    `json:"has_problems"`
	Official_line    string  `json:"official_line"`
	Delay            *int    `json:"delay"`
}

func (r ServiceRaw) ToDomain() Service {
	return Service{
		LineInfo: LineInfo{
			Line:             r.Service,
			Destination:      r.Destination,
			Journey_code:     r.Codice_corsa,
		},
		Arrival_time:     r.Arrival,
		State:            r.Type,
		Basin:            r.FleetCode,
		Vehicle_table:    r.DutyId,
		Vehicle:          r.Busnum,
		Line_type:        r.ServiceType,
		Occupancy_status: r.OccupancyStatus,
		Total_room:       r.Posti_totali,
		Passenger_number: r.Num_passeggeri,
		Next_stop:        r.Next_stop,
	}
}

type ArrivalRaw struct {
	Arrival struct {
		Services []ServiceRaw `json:"services"`
	} `json:"arrival"`
}

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