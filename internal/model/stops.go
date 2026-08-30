package model

import "time"

type Stop struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome" json:"name"`
}

type StopCorrection struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome_mod" json:"name"`
}

type StopsInfo struct {
	Count      int       `json:"count"`
	Updated_at time.Time `json:"updated_at"`
}
