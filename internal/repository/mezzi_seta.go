package repository

import (
	"database/sql"
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

func GetVehicleLastSeen(id string) (*time.Time, error) {
	var result *time.Time
	err := DB_MEZZI.Get(&result, "SELECT last_seen FROM mezzi_seta WHERE matricola = ?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("[GetVehicleLastSeen] errore di lettura db:", err)
		} else {
			return nil, err
		}
	}

	return result , nil
}

func GetVehiclePlate(id string) (*string, error) {
	var result *string
	err := DB_MEZZI.Get(&result, "SELECT targa FROM mezzi_seta WHERE matricola = ?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("[GetVehicleLastSeen] errore di lettura db:", err)
		} else {
			return nil, err
		}
	}

	return result , nil
}

func UpdateVehicleStatus(id string, status string, last_seen time.Time) {
	_, err := DB_MEZZI.Exec("UPDATE mezzi_seta SET stato = ?, last_seen = ? WHERE matricola = ?", status, last_seen, id)
	if err != nil {
		fmt.Println("[UpdateVehicleStatus] error writing to db:", err)
	}
}
