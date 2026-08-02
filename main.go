package main

import (
	"fmt"
	"log"
	"net/http"

	"setaapi/internal/handler"
)

func main() {
	mux := http.NewServeMux();
	mux.HandleFunc("GET /health", handler.HealthCheckHandler);
	mux.HandleFunc("GET /arrivals/{id}", handler.ArrivalsHandler);

	//Listen on port and start API
	fmt.Println("Server avviato");
	log.Print(http.ListenAndServe(":5001",mux));
}