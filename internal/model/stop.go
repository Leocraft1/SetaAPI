package model

type Stop struct {
	Codice string 	`db:"codice" json:"codice"`
	Nome string 	`db:"nome" json:"nome"`
}