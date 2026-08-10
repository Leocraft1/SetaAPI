package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetRoute(code string) model.LineeRC{
	var result model.LineeRC
	err := DB_CONTENT.Get(&result, "SELECT * FROM linee_rc WHERE rc = ?", code)
	if err != nil {
		fmt.Println("[GetRoute] Errore di lettura db:", err)
	}
	return result
}

func GetRoutesDistinct() []string {
	var result []string
	err := DB_CONTENT.Select(&result, "SELECT DISTINCT linea FROM linee_rc")
	if err != nil {
		fmt.Println("[GetRoutesDistinct] Errore di lettura db:", err)
	}
	return result
}

func GetRCTable() []model.LineeRC {
	var result []model.LineeRC
	err := DB_CONTENT.Select(&result, "SELECT * FROM linee_rc")
	if err != nil {
		fmt.Println("[GetRouteCodes] Errore di lettura db:", err)
	}
	return result
}