package repository

import (
	"fmt"
	"regexp"
	"setaapi/internal/model"
	"sort"
	"strconv"
	"unicode"
)

func GetAssignments() []model.Assignment {
	var results []model.Assignment
	err := DB_CONTENT.Select(&results, "SELECT * FROM assignments")
	if err != nil {
		fmt.Println("GetAssignments db error:", err)
	}

	return sortByTable(results)
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


// ---------------------
// - PRIVATE FUNCTIONS -
// ---------------------
//Sorts assignments
func sortByTable(buses []model.Assignment) []model.Assignment {
	sort.SliceStable(buses, func(i, j int) bool {
		numI := extractLineNumber(buses[i].VehicleTable)
		numJ := extractLineNumber(buses[j].VehicleTable)
		if numI != numJ {
			return numI < numJ
		}
		return buses[i].VehicleTable < buses[j].VehicleTable
	})

	numeric := make([]model.Assignment, 0, len(buses))
	letters := make([]model.Assignment, 0)

	for _, b := range buses {
		if len(b.VehicleTable) > 0 && unicode.IsLetter(rune(b.VehicleTable[0])) {
			letters = append(letters, b)
		} else {
			numeric = append(numeric, b)
		}
	}

	return append(numeric, letters...)
}

var numericPartRegex = regexp.MustCompile(`\d+`)

func extractLineNumber(line string) int {
	match := numericPartRegex.FindString(line)
	if match == "" {
		return 0
	}
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return num
}
