package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"setaapi/internal/model"
)

type query1Result struct {
	Model *string `db:"modello"`
	Ramp  *int    `db:"pedana"`
}

func GetModelRampByCode(code string) (*string, *int) {
	var m query1Result
	err := DB_MEZZI.Get(&m, "SELECT COALESCE(NULLIF(short_model, ''), modello) AS modello, pedana FROM mezzi_seta WHERE matricola = ?", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//Do nothing
		} else {
			fmt.Println("[GetModelRampByCode] Errore di lettura db:", err)
		}
	}
	return m.Model, m.Ramp
}

func GetModelsDistinct() []string {
	var result []string
	err := DB_MEZZI.Select(&result, "SELECT DISTINCT COALESCE(NULLIF(short_model, ''), modello) FROM mezzi_seta ORDER BY COALESCE(NULLIF(short_model, ''), modello)")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//Do nothing
		} else {
			fmt.Println("[GetModelRampByCode] Errore di lettura db:", err)
		}
	}
	return result
}

func GetAEP() []model.HasAEP {
	var result []model.HasAEP
	err := DB_MEZZI.Select(&result, "SELECT matricola, has_aep FROM mezzi_seta")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//Do nothing
		} else {
			fmt.Println("[GetModelRampByCode] Errore di lettura db:", err)
		}
	}

	return result
}
