package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetUpdateTimestamps() []model.UpdateEntry {
	var results []model.UpdateEntry
	err := DB_CONTENT.Select(&results, "SELECT * FROM update_timestamps")
	if err != nil {
		fmt.Println("[GetStopCount] db error:", err)
	}

	return results
}