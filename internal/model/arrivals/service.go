package arrivals

type Service struct {
	Line string `json:"line"`;
	Arrival_time string `json:"arrival_time"`;
	State string `json:"state"`;
	Destination string `json:"destination"`;
	Basin string `json:"basin"`;
	Vehicle_table string `json:"vehicle_table"`;
	Vehicle string `json:"vehicle"`;
	Line_type string `json:"line_type"`;
	Occupancy_status *string `json:"occupancy_status"`;
	Route_code string `json:"route_code"`
	Total_room *int `json:"total_room"`
	Passenger_number *int `json:"passenger_number"`
	Next_stop *string `json:"next_stop"`
	Has_problems bool `json:"has_problems"`
	Official_line string `json:"official_line"`
	Delay *int `json:"delay"`
}

func (r ServiceRaw) ToDomain() Service {
    return Service{
        Line:  r.Service,
        Arrival_time: r.Arrival,
        State: r.Type,
		Destination: r.Destination,
		Basin: r.FleetCode,
		Vehicle_table: r.DutyId,
		Vehicle: r.Busnum,
		Line_type: r.ServiceType,
		Occupancy_status: r.OccupancyStatus,
		Route_code: r.Codice_corsa,
		Total_room: r.Posti_totali,
		Passenger_number: r.Num_passeggeri,
		Next_stop: r.Next_stop,
    }
}