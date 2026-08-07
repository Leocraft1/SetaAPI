package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type queryResult struct {
	Model *string `db:"modello"`
	Ramp  *int    `db:"pedana"`
}

func GetModelRampByCode(code string) (*string, *int) {
	var m queryResult
	err := DB_MEZZI.Get(&m, "SELECT modello, pedana FROM mezzi_seta WHERE matricola = ?", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Nessuna riga trovata per matricola ", code)
		} else {
			fmt.Println("Errore di lettura db: ", err)
		}
	}
	return m.Model, m.Ramp
}
