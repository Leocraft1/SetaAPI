package model

type Stop struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome" json:"name"`
}

type StopCorrection struct {
	Code string `db:"codice" json:"code"`
	Name string `db:"nome_mod" json:"name"`
}

type StopsInfo struct {
	Count           int    `json:"count"`
	Updated_at_date string `json:"updated_at_date"`
	Updated_at_time string `json:"updated_at_time"`
}
