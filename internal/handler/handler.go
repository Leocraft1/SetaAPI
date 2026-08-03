package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"setaapi/internal/model/arrivals"
	"setaapi/internal/service"
)

//URLs decl. section
var arrivalsBaseUrl string = "https://avm.setaweb.it/SETA_WS/services/arrival/"

//GET /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	out := "API is healthy,";

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(out);
}

//GET /arrivals/{id}
func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	stopId := r.PathValue("id");
	response, err := http.Get(arrivalsBaseUrl+stopId);

	//Checks for connection errors and (TODO) HTTP code
	if err != nil {
		fmt.Println("ArrivalHandler error: ", err);
	}

	//Parses response into struct
	var arrivalsRaw arrivals.ArrivalRaw;
	json.NewDecoder(response.Body).Decode(&arrivalsRaw);
	
	//Fixes incorrect stuff/add parameters
	arrivals := service.FixArrivals(arrivalsRaw);

	//Set headers
	w.Header().Set("Content-Type", "application/json");
	w.WriteHeader(response.StatusCode);

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(arrivals);
}