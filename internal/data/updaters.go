package data

import (
	"fmt"
	"setaapi/internal/handler"
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"setaapi/internal/service"
)

func UpdateStops() {
	buses, err := service.GetBusesinservice(handler.WimbBaseUrl, handler.LineeDynUrl)
	if err != nil {
		fmt.Println("UpdateStops error getting busesinservice response:", err)
		return
	}
	var currStops []model.Stop
	for _, val := range buses.Buses {
		if val.Next_stop != "" {
			newCurrStop := model.Stop{
				Code: val.Stop_code,
				Name: val.Next_stop,
			}
			currStops = append(currStops, newCurrStop)
		}
	}

	repository.SaveStops(currStops)
	fmt.Println("UpdateStops OK")
}

func UpdateRoutes() {
	buses, err := service.GetBusesinservice(handler.WimbBaseUrl, handler.LineeDynUrl)
	if err != nil {
		fmt.Println("UpdateRoutes error getting busesinservice response:", err)
		return
	}

	var currRoutes []model.Route
	for _, val := range buses.Buses {
		if val.Next_stop != "" {
			newCurrRoute := model.Route{
				Linea: val.Official_line,
				Rc: val.Route_code,
				Still_exists: true,
			}
			currRoutes = append(currRoutes, newCurrRoute)
		}
	}

	repository.SaveRoutes(currRoutes)
	fmt.Println("UpdateRoutes OK")
}