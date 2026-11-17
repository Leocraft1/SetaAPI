package data

import (
	"fmt"
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"setaapi/internal/repository/static"
	"time"
)

// Checks if vehicles are servicing (have been assigned), if so sets last seen date and resets status if present.
func SyncVehicleStatus() {
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	assignments := static.GetAssignments()
	for _, val := range assignments {
		lastSeen, err := repository.GetVehicleLastSeen(val.Vehicle)
		if repository.CheckKnownVehicle(val.Vehicle) && (lastSeen == nil || !lastSeen.Equal(today)) && err == nil {
			//Sets "" state (active vehicle)
			repository.UpdateVehicleStatus(val.Vehicle, "", today)
		}
	}

	fmt.Println("SyncVehicleStatus OK")
}

func SyncAssignments() {
	var vehicles []model.BusInService

	vehiclesMap := make(map[string]string)
	for _, val := range vehicles {
		vehiclesMap["val.Vehicle_table"] = val.Vehicle
	}

	old := make(map[string]string)
	assMap := static.GetAssignmentsMap()
	//Searches for stuff that doesn't communicate GPS anymore
	for idx, val := range assMap {
		_, ok := vehiclesMap[idx]
		if !ok {
			old[idx] = val
		}
	}
	static.InsertGPSAssignments(vehiclesMap)
	static.SetNoGPS(old)
	fmt.Println("SyncAssignments OK")
}
