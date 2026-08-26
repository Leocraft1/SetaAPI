package repository

import (
	"fmt"
	"net/http"
	"setaapi/internal/model"
)

// URLs decl. section
var RoutestopsBaseUrl = "https://wimb.setaweb.it/publicmapbe/waypoints/getroutewaypoints/"

func GetRoute(code string) model.LineeRC {
	var result model.LineeRC
	err := DB_CONTENT.Get(&result, "SELECT * FROM routes WHERE rc = ?", code)
	if err != nil {
		fmt.Println("[GetRoute] Errore di lettura db:", err)
	}
	return result
}

func GetLinesDistinct() []string {
	var result []string
	err := DB_CONTENT.Select(&result, "SELECT DISTINCT linea FROM routes")
	if err != nil {
		fmt.Println("[GetLinesDistinct] Errore di lettura db:", err)
	}
	return result
}

func GetRCTable() []model.LineeRC {
	var result []model.LineeRC
	err := DB_CONTENT.Select(&result, "SELECT * FROM routes")
	if err != nil {
		fmt.Println("[GetRCTable] Errore di lettura db:", err)
	}
	return result
}

func GetExists() []model.StillExists {
	var results []model.StillExists
	err := DB_CONTENT.Select(&results, "SELECT rc, still_exists FROM routes")
	if err != nil {
		fmt.Println("[GetExists] Errore di lettura db:", err)
	}

	return results
}

// Updates route status in route table (still_exists column)
func UpdateRoutesStatus() {
	routeCodes := GetExists()

	rcMap := make(map[string]bool)
	for _, val := range routeCodes {
		rcMap[val.Rc] = val.Still_exists
	}

	newStatusMap := make(map[string]bool)
	for idx, val := range rcMap {
		response, err := http.Get(RoutestopsBaseUrl + idx)
		if err != nil {
			fmt.Println("UpdateRoutesStatus error connecting to upstream:", err)
		}

		if response.StatusCode == 404 && val {
			newStatusMap[idx] = false
		} else if response.StatusCode == 200 && !val {
			newStatusMap[idx] = true
		}
	}

	for idx, val := range newStatusMap {
		_, err := DB_CONTENT.Exec("UPDATE routes SET still_exists = ? WHERE rc = ?", val, idx)
		if err != nil {
			fmt.Println("UpdateRoutesStatus db error:", err)
		}
	}
}
