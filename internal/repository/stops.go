package repository

import (
	"fmt"
	"setaapi/internal/model"
	"time"
)

func GetStops() []model.Stop {
	var result []model.Stop
	err := DB_CONTENT.Select(&result, "SELECT * FROM stops")
	if err != nil {
		fmt.Println("[GetStops] db error:", err)
	}

	return result
}

func GetStopCount() (int, time.Time) {
	var rowCount []int
	var updateTS []time.Time
	err := DB_CONTENT.Select(&rowCount, "SELECT COUNT(*) FROM stops")
	err = DB_CONTENT.Select(&updateTS, "SELECT updated_at FROM update_timestamps WHERE table_name = ?", "stops")
	if err != nil {
		fmt.Println("[GetStopCount] db error:", err)
	}

	rows := rowCount[0]
	updateTime := updateTS[0]
	return rows, updateTime
}

func SaveStops(stops []model.Stop) {
	dbData := GetStops()

	dbMap := make(map[string]bool)
	for _, val := range dbData {
		dbMap[val.Code] = true
	}

	var new []model.Stop
	for _, val := range stops {
		_, ok := dbMap[val.Code]
		if !ok {
			newStop := model.Stop{
				Code: val.Code,
				Name: val.Name,
			}
			new = append(new, newStop)
		}
	}

	//Corrects stop information if wrong
	corrections := GetStopCorrections()
	for idx, val := range new {
		val2, ok := corrections[val.Code]
		if ok {
			new[idx].Name = val2
		}
	}

	//Database insert
	for _, val := range new {
		_, err := DB_CONTENT.Exec("INSERT INTO stops VALUES(?, ?)", val.Code, val.Name)
		if err != nil {
			fmt.Println("[SaveStops] db error:", err)
		}
	}

	//Updates timestamp in update_timestamps
	if len(new) > 0 {
		_, err := DB_CONTENT.Exec("UPDATE update_timestamps SET updated_at = ? WHERE table_name = ?", time.Now(), "stops")
		if err != nil {
			fmt.Println("[SaveStops] db error updating timestamp:", err)
		}
	}
}
