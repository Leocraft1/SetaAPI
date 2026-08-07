package busesinservice

import "setaapi/internal/model"

type Bus struct {
	model.LineInfo
	Line_type         string  `json:"line_type"`
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
		LineInfo: model.LineInfo{
			Line: r.Properties.Linea,
			Destination: r.Properties.Route_desc,
			Journey_code: r.Properties.Journey_code,
		},
		Line_type: r.Properties.Service_tag,
		Delay: r.Properties.Delay,
		Vehicle_table: r.Properties.Duty_id,
		Vehicle: r.Properties.Vehicle_code,
		Model: r.Properties.Model,
		Plate_num: r.Properties.Plate_num,
		Ramp: r.Properties.Pedana,
		Next_stop: r.Properties.Next_stop,
		Stop_code: r.Properties.Waypoint_code,
		Route_code: r.Properties.Route_code,
		Total_room: r.Properties.Posti_totali,
		Occupancy_lastupd: r.Properties.Occupancy_lastupd,
		Passenger_number: r.Properties.Num_passeggeri,
	}
}