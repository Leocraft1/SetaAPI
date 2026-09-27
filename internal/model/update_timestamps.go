package model

import "time"

type UpdateEntry struct {
	TableName    string    `db:"table_name" json:"table_name"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
	UpdatedAtStr string `json:"updated_at"`
}
