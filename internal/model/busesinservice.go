package model

type Buses struct {
	Buses []Bus `json:"buses"`
}

type Bus struct {
	LineInfo
	Delay             int     `json:"delay"`
	Vehicle_table     string  `json:"vehicle_table"`
	Vehicle           string  `json:"vehicle"`
	Model             string  `json:"model"`
	Plate_num         string  `json:"plate_num"`
	Ramp              int     `json:"ramp"`
	Next_stop         string  `json:"next_stop"`
	Stop_code         string  `json:"stop_code"`
	Route_code        string  `json:"route_code"`
	Total_room        int     `json:"total_room"`
	Occupancy_lastupd *string `json:"occupancy_lastupd"`
	Passenger_number  *int    `json:"passenger_number"`
	Has_problems      bool    `json:"has_problems"`
	Has_AEP           bool    `json:"has_AEP"`
}

func (r BusRaw) ToDomain() Bus {
	return Bus{
		LineInfo: LineInfo{
			Line:         r.Properties.Linea,
			Destination:  r.Properties.Route_desc,
			Journey_code: r.Properties.Journey_code,
			Line_type:    r.Properties.Service_tag,
		},
		Delay:             r.Properties.Delay,
		Vehicle_table:     r.Properties.Duty_id,
		Vehicle:           r.Properties.Vehicle_code,
		Model:             r.Properties.Model,
		Plate_num:         r.Properties.Plate_num,
		Ramp:              r.Properties.Pedana,
		Next_stop:         r.Properties.Next_stop,
		Stop_code:         r.Properties.Waypoint_code,
		Route_code:        r.Properties.Route_code,
		Total_room:        r.Properties.Posti_totali,
		Occupancy_lastupd: r.Properties.Occupancy_lastupd,
		Passenger_number:  r.Properties.Num_passeggeri,
	}
}

type BusesRaw struct {
	Properties []BusRaw `json:"features"`
}

type BusRaw struct {
	Properties struct {
		Linea             string  `json:"linea"`
		Route_desc        string  `json:"route_desc"`
		Service_tag       string  `json:"service_tag"`
		Delay             int     `json:"delay"`
		Duty_id           string  `json:"duty_id"`
		Vehicle_code      string  `json:"vehicle_code"`
		Model             string  `json:"model"`
		Plate_num         string  `json:"plate_num"`
		Pedana            int     `json:"pedana"`
		Next_stop         string  `json:"next_stop"`
		Waypoint_code     string  `json:"waypoint_code"`
		Route_code        string  `json:"route_code"`
		Journey_code      string  `json:"journey_code"`
		Posti_totali      int     `json:"posti_totali"`
		Occupancy_lastupd *string `json:"occupancy_lastupd"`
		Num_passeggeri    *int    `json:"num_passeggeri"`
	}
}
