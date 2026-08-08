package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"setaapi/config"
	"setaapi/internal/model/arrivals"
	"setaapi/internal/model/busesinservice"
	"setaapi/internal/service"
)

// URLs decl. section
var arrivalsBaseUrl string = "https://avm.setaweb.it/SETA_WS/services/arrival/"
var wimbBaseUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO"

// GET /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	message := "API is healthy and running."

	w.Write([]byte(message))
}

// GET /arrivals/{id}
func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	stopId := r.PathValue("id")
	response, err := http.Get(arrivalsBaseUrl + stopId)

	//Checks for connection errors and (TODO) HTTP code
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
	}

	//Parses response into struct
	var arrivalsRaw arrivals.ArrivalRaw
	json.NewDecoder(response.Body).Decode(&arrivalsRaw)

	//Fixes incorrect stuff/add parameters
	arrivals := service.FixArrivals(arrivalsRaw)

	//Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(arrivals)
}

// GET /busesinservice
func BusesinserviceHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(wimbBaseUrl)
	if err != nil {
		fmt.Println("BusesinserviceHandler error: ", err)
		w.Write([]byte("BusesinserviceHandler error: " + err.Error()))
	}
	var busesRaw busesinservice.BusesRaw
	json.NewDecoder(response.Body).Decode(&busesRaw)

	buses := service.FixBusesinservice(busesRaw)

	//Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(buses)
}

// GET /vehicleinfo/{id}
// Searches for requested vehicle on /busesinservice response
func VehicleinfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	response, err := http.Get("http://localhost"+config.PORT+"/busesinservice")
	if err != nil {
		fmt.Println("VehicleinfoHandler error: ", err)
		w.Write([]byte("VehicleinfoHandler error: " + err.Error()))
	}
	var buses busesinservice.Buses
	json.NewDecoder(response.Body).Decode(&buses)

	for idx, val := range buses.Buses {
		if val.Vehicle == id {
			//Set headers
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.StatusCode)

			//Parses response back to JSON and returns it to client
			json.NewEncoder(w).Encode(buses.Buses[idx])
			return
		}
	}
	//Vehicle not found
	w.WriteHeader(404)
	w.Write([]byte("Vehicle not found."))
}
