package handler

import (
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request){
	//out := []string {"[Ciao]"};
	out := "API is running."

	//w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out);
}