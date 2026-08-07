package main

import (
	"fmt"
	"log"
	"net/http"
	"setaapi/internal/handler"
	"setaapi/internal/repository"
)

func main() {
	//DBs Init
	repository.InitMezzi()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	mux.HandleFunc("GET /arrivals/{id}", handler.ArrivalsHandler)
	mux.HandleFunc("GET /busesinservice", handler.BusesinserviceHandler)

	//Listen on port and start API
	fmt.Println("Server avviato")
	log.Print(http.ListenAndServe(":5001", mux))
}
