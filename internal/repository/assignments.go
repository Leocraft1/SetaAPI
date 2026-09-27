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
	//Insert (gps and last_seen are true)
	for idx, val := range new {
		_, err := DB_CONTENT.Exec("INSERT INTO assignments VALUES(?, ?, ?, ?)", idx, val, true, true)
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