package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetStopCorrections() map[string]string {
	var results []model.StopCorrection
	err := DB_CONTENT.Select(&results, "SELECT * FROM stops_correction")
	if err != nil {
		fmt.Println("GetStopCorrections db error:", err)
	}

	resultsMap := make(map[string]string)
	for _, val := range results {
		resultsMap[val.Code] = val.Name
	}
	
	return resultsMap
}