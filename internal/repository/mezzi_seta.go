package repository

import (
	"fmt"
	"time"
)

func CheckKnownVehicle(id string) bool {
	var result int
	err := DB_MEZZI.Get(&result, "SELECT COUNT(*) FROM mezzi_seta WHERE matricola = ?", id)
	if err != nil {
		fmt.Println("[CheckVehicleExists] errore di lettura db:", err)
	}
	if result == 0 {
		return false
	} else {
		return true
	}
}

func GetVehicleLastSeen(id string) *time.Time {
	var result *time.Time
	err := DB_MEZZI.Get(&result, "SELECT last_seen FROM mezzi_seta WHERE matricola = ?", id)
	if err != nil {
		fmt.Println("[GetVehicleLastSeen] errore di lettura db:", err)
	}

	//Vehicle should exist when this func is called
	return result
}

func UpdateVehicleStatus(id string, status string, last_seen time.Time) {
	_, err := DB_MEZZI.Exec("UPDATE mezzi_seta SET stato = ?, last_seen = ? WHERE matricola = ?", status, last_seen, id)
	if err != nil {
		fmt.Println("[UpdateVehicleStatus] error writing to db:", err)
	}
}