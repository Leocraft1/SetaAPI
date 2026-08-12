package model

type RouteStops struct {
	Still_exists bool `json:"still_exists"`
	Stops []RouteStop `json:"stops"`
}

type RouteStop struct {
	Code string `db:"codice_fermata" json:"code"`
	Name string `db:"nome_fermata" json:"name"`
	Is_last bool `db:"is_last" json:"is_last"`
}