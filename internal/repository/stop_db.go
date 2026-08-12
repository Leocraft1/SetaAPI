package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetStops() []model.Stop {
	var result []model.Stop
	err := DB_CONTENT.Select(&result, "SELECT * FROM fermate")
	if err != nil {
		fmt.Println("[GetStops] Errore di lettura db:", err)
	}
	return result
}
