package data

import (
	"fmt"
	"setaapi/internal/handler"
	"setaapi/internal/repository"
	"setaapi/internal/service"
	"time"
)

//Checks if vehicles have gone out, if so sets last seen date and resets status if present
func SyncVehicleStatus() {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	vehicles, err := service.GetBusesinservice(handler.WimbBaseUrl, handler.LineeDynUrl)
	if err != nil {
		fmt.Println("SyncVehicleStatus unable to get buses in service:", err)
	}
	
	for _, val := range vehicles.Buses {
		lastSeen := repository.GetVehicleLastSeen(val.Vehicle)
		if repository.CheckKnownVehicle(val.Vehicle) && (lastSeen == nil || !lastSeen.Equal(today)) {
			//Sets "" state (active vehicle)
			repository.UpdateVehicleStatus(val.Vehicle, "", today)
		}
	}

	fmt.Println("SyncVehicleStatus OK")
}