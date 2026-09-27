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

func SyncAssignments() {
	vehicles, err := service.GetBusesinservice(handler.WimbBaseUrl, handler.LineeDynUrl)
	if err != nil {
		fmt.Println("SyncAssignments unable to get buses in service:", err)
	}

	vehiclesMap := make(map[string]string)
	for _, val := range vehicles.Buses {
		vehiclesMap[val.Vehicle_table] = val.Vehicle
	}

	new := make(map[string]string)
	assMap := repository.GetAssignmentsMap()
	for idx, val := range vehiclesMap {
		_, ok := assMap[idx]
		if !ok {
			new[idx] = val
		}
	}
	repository.InsertGPSAssignments(new)
	fmt.Println("SyncAssignments OK")
}