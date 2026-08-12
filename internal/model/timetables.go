package model

type Timetable struct {
	DataPercorso string     `json:"data_percorso"`
	Stops        []string   `json:"stops"`
	Timetable    [][]string `json:"timetable"`
}

type TimetableResponse struct {
	Stops    []string  `json:"stops"`
	Journeys []Journey `json:"journeys"`
}

type Journey struct {
	RouteCode string   `json:"route_code"`
	Disp_line *string  `json:"display_line"`
	Disp_dest *string  `json:"display_destination"`
	Times     []string `json:"times"`
}
