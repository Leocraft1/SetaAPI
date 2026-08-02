package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	out := "API is healthy,";

	w.Header().Set("Content-Type", "application/json");
	json.NewEncoder(w).Encode(out);
}

func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	stopId := r.PathValue("id");
	response, err := http.Get(`https://avm.setaweb.it/SETA_WS/services/arrival/`+stopId);

	if err != nil {
		fmt.Println("ArrivalHandler error: ", err);
	}

	w.Header().Set("Content-Type", "application/json");
	//json.NewEncoder(w).Encode(response.Body);
	bytes, _ := io.ReadAll(response.Body)
	w.Write(bytes)
}