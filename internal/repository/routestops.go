package repository

import (
	"fmt"
	"setaapi/internal/model"
)

func GetRouteStops(rc string) model.RouteStops {
	var result []model.RouteStop
	err := DB_CONTENT.Select(&result, "SELECT nome_fermata, codice_fermata, is_last FROM routestops WHERE rc = ? ORDER BY `order`", rc)
	if err != nil {
		fmt.Println("[GetRouteByRCStops] Errore di lettura db:", err)
	}
	var results model.RouteStops
	route := GetRouteByRC(rc)
	results.Still_exists = route.Still_exists
	results.Stops = result
	return results
}