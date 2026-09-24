package repository

import (
	"fmt"
	"io"
	"net/http"
	"setaapi/internal/model"
	"time"
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
		_, err := DB_CONTENT.Exec("INSERT INTO routes VALUES(?, ?, ?, ?, ?, ?)", val.Linea,val.Rc, val.Disp_linea, val.Disp_dest, val.Desc, val.Still_exists)
		if err != nil {
			fmt.Println("[SaveRoutes] db error:", err)
		}
	}

	//Updates timestamp in update_timestamps
	if len(new) > 0 {
		_, err := DB_CONTENT.Exec("UPDATE update_timestamps SET updated_at = ? WHERE table_name = ?", time.Now(), "routes")
		if err != nil {
			fmt.Println("[SaveStops] db error updating timestamp:", err)
		}
	}
}

// Client condiviso con timeout: http.Get usa DefaultClient, che non ha timeout
var httpClient = &http.Client{Timeout: 10 * time.Second}

// Updates route status in routes table (still_exists column)
func UpdateRoutesStatus() {
	routeCodes := GetExists()

	rcMap := make(map[string]bool, len(routeCodes))
	for _, val := range routeCodes {
		rcMap[val.Rc] = val.Still_exists
	}

	newStatusMap := make(map[string]bool)
	for rc, exists := range rcMap {
		status, err := fetchStatus(RoutestopsBaseUrl + rc)
		if err != nil {
			fmt.Println("UpdateRoutesStatus error connecting to upstream:", err)
			continue
		}

		if status == http.StatusNotFound && exists {
			newStatusMap[rc] = false
		} else if status == http.StatusOK && !exists {
			newStatusMap[rc] = true
		}
	}

	for rc, val := range newStatusMap {
		if _, err := DB_CONTENT.Exec("UPDATE routes SET still_exists = ? WHERE rc = ?", val, rc); err != nil {
			fmt.Println("UpdateRoutesStatus db error:", err)
		}
	}
}

func fetchStatus(url string) (int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// Svuota il body così la connessione keep-alive può essere riutilizzata
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}