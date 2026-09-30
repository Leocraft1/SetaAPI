package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetAssignments() []model.Assignment {
	var results []model.Assignment
	err := DB_CONTENT.Select(&results, "SELECT * FROM assignments")
	if err != nil {
		fmt.Println("GetAssignments db error:", err)
	}

	return results
}

func GetAssignmentsMap() map[string]string {
	var results []model.Assignment
	err := DB_CONTENT.Select(&results, "SELECT * FROM assignments")
	if err != nil {
		fmt.Println("GetAssignments db error:", err)
	}

	resultsMap := make(map[string]string)
	for _, val := range results {
		resultsMap[val.VehicleTable] = val.Vehicle
	}
	
	return resultsMap
}

func InsertGPSAssignments(new map[string]string) {
	//Insert (gps and last_seen are true) overwriting if vehicle was already set
	for idx, val := range new {
		_, err := DB_CONTENT.Exec(
			`INSERT INTO assignments (vehicle_table, vehicle, is_GPS)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE vehicle = VALUES(vehicle), is_GPS = VALUES(is_GPS), wants_last_seen = VALUES(wants_last_seen)`,
			idx, val, true,
		)
		if err != nil {
			fmt.Println("InsertGPSAssignments db error:", err)
		}
	}
}

func SetNoGPS(old map[string]string) {
	for idx := range old {
		_, err := DB_CONTENT.Exec("UPDATE assignments SET is_GPS = false WHERE vehicle_table = ?", idx)
		if err != nil {
			fmt.Println("InsertGPSAssignments db error:", err)
		}
	}
}

func DeleteAllAssignments() {
	_, err := DB_CONTENT.Exec("DELETE FROM assignments")
	if err != nil {
		fmt.Println("DeleteAllAssignments db error:", err)
	}
}