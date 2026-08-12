package service

import (
	"setaapi/internal/model"
	"setaapi/internal/repository"
)

func GetRoutestops(route_code string) model.RouteStops{
	results := repository.GetRouteStops(route_code)

	return results
}