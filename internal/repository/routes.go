package repository

import (
	"fmt"
	"net/http"
	"setaapi/internal/model"
)

// URLs decl. section
var RoutestopsBaseUrl = "https://wimb.setaweb.it/publicmapbe/waypoints/GetRouteByRCwaypoints/"

func GetRouteByRC(code string) model.Route {
	var result model.Route
	err := DB_CONTENT.Get(&result, "SELECT * FROM routes WHERE rc = ?", code)
	if err != nil {
		fmt.Println("[GetRouteByRC] Errore di lettura db:", err)
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

func GetRoutes() []model.Route {
	var result []model.Route
	err := DB_CONTENT.Select(&result, "SELECT * FROM routes")
	if err != nil {
		fmt.Println("[GetRoutes] Errore di lettura db:", err)
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

func SaveRoutes(routes []model.Route) {
	dbData := GetRoutes()

	dbMap := make(map[string]bool)
	for _, val := range dbData {
		dbMap[val.Rc] = true
	}

	var new []model.Route
	for _, val := range routes {
		_, ok := dbMap[val.Rc]
		if !ok {
			newRoute := model.Route{
				Linea: val.Linea,
				Rc: val.Rc,
				Disp_linea: val.Disp_linea,
				Disp_dest: val.Disp_dest,
				Desc: val.Desc,
				Still_exists: val.Still_exists,
			}

			new = append(new, newRoute)
		}
	}

	//Database insert
	for _, val := range new {
		_, err := DB_CONTENT.Exec("INSERT INTO stops VALUES(?, ?, ?, ?, ?, ?)", val.Linea,val.Rc, val.Disp_linea, val.Disp_dest, val.Desc, val.Still_exists)
		if err != nil {
			fmt.Println("[SaveStops] db error:", err)
		}
	}
}

// Updates route status in routes table (still_exists column)
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
