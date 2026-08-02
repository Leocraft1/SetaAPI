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

	//Listen on port and start API
	fmt.Printf("Server avviato");
	log.Print(http.ListenAndServe(":5001",mux));
}