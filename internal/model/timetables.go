package model

type TimetableResponse struct {
	Timetable []Timetable `json:"response"`
}

type Timetable struct {
	DataPercorso string     `json:"data_percorso"`
	Stops        []string   `json:"stops"`
	Timetable    [][]string `json:"timetable"`
}
