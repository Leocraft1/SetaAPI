package model

type Assignment struct {
	VehicleTable string `db:"vehicle_table" json:"vehicle_table"`
	Vehicle string `db:"vehicle" json:"vehicle"`
	IsGPS bool `db:"is_GPS" json:"is_GPS"`
	WantsLastSeen bool `db:"wants_last_seen" json:"wants_last_seen"`
}