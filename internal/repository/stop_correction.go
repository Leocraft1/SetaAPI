package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetStopCorrections() []model.StopCorrection {
	var results []model.StopCorrection
	err := DB_CONTENT.Select(&results, "SELECT * FROM stops_correction")
	if err != nil {
		fmt.Println("GetStopCorrections db error:", err)
	}
	
	return results
}