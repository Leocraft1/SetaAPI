package model

type Stop struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome" json:"name"`
}

type StopCorrection struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome_mod" json:"name"`
}