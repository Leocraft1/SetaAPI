package model

type Assignment struct {
	VehicleTable string `db:"vehicle_table" json:"vehicle_table"`
	Vehicle      string `db:"vehicle" json:"vehicle"`
	IsGPS        bool   `db:"is_GPS" json:"is_GPS"`
}

type CreateAssignmentRequest struct {
	VehicleTable string `json:"vehicle_table"`
	Vehicle      int    `json:"vehicle"`
}

type DeleteAssignmentRequest struct {
	VehicleTable string `json:"vehicle_table"`
}