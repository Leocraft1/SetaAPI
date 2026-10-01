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
	//Insert (gps is true) overwriting if vehicle was already set
	for idx, val := range new {
		_, err := DB_CONTENT.Exec(
			`INSERT INTO assignments (vehicle_table, vehicle, is_GPS)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE vehicle = VALUES(vehicle), is_GPS = VALUES(is_GPS)`,
			idx, val, true,
		)
		if err != nil {
			fmt.Println("InsertGPSAssignments db error:", err)
		}
	}
}

func InsertAssignment(table string, vehicle int) error {
	//Insert (gps is false) overwriting if vehicle was already set
	_, err := DB_CONTENT.Exec(
		`INSERT INTO assignments (vehicle_table, vehicle, is_GPS)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE vehicle = VALUES(vehicle), is_GPS = VALUES(is_GPS)`,
		table, vehicle, false,
	)

	return err
}

func UpdateAssignment(table string, vehicle int) error {
	//New vehicle for given table
	_, err := DB_CONTENT.Exec(
		`UPDATE assignments SET vehicle = ? WHERE vehicle_table = ?`,
		vehicle, table,
	)
	
	return err
}

func DeleteAssignment(table string) error {
	//Delete assignment
	_, err := DB_CONTENT.Exec(
		`DELETE FROM assignments WHERE vehicle_table = ?`,
		table,
	)
	
	return err
}

func SetNoGPS(old map[string]string) {
	for idx := range old {
		_, err := DB_CONTENT.Exec("UPDATE assignments SET is_GPS = false WHERE vehicle_table = ?", idx)
		if err != nil {
			fmt.Println("SetNoGPS db error:", err)
		}
	}
}

func DeleteAllAssignments() {
	_, err := DB_CONTENT.Exec("DELETE FROM assignments")
	if err != nil {
		fmt.Println("DeleteAllAssignments db error:", err)
	}
}