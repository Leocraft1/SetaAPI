package busesinservice

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