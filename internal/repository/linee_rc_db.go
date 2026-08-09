package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetRoute(code string) model.LineeRC{
	var result model.LineeRC
	err := DB_CONTENT.Get(&result, "SELECT * FROM linee_rc WHERE rc = ?", code)
	if err != nil {
		fmt.Println("GetRoute DB error:", err)
	}
	return result
}