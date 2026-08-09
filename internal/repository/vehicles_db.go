package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
)

type query1Result struct {
	Model *string `db:"modello"`
	Ramp  *int    `db:"pedana"`
}

func GetModelRampByCode(code string) (*string, *int) {
	var m query1Result
	err := DB_MEZZI.Get(&m, "SELECT modello, pedana FROM mezzi_seta WHERE matricola = ?", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("[GetModelRampByCode] Nessuna riga trovata per matricola:", code)
		} else {
			fmt.Println("[GetModelRampByCode] Errore di lettura db:", err)
		}
	}
	return m.Model, m.Ramp
}

func GetModelsDistinct() []string {
	var result []string
	err := DB_MEZZI.Select(&result, "SELECT DISTINCT modello FROM mezzi_seta")
	if err != nil {
		fmt.Println("[GetModelsDistinct] Errore di lettura db: ", err)
	}
	//Sorts alphabetically
	slices.Sort(result)
	return result
}
