package main

import (
	"fmt"
	"log"
	"net/http"
	"setaapi/config"
	"setaapi/internal/handler"
	"setaapi/internal/repository"
)

func main() {
	//Init config
	config.LoadConf()
	//DBs Init
	repository.InitMezzi()
	repository.InitContent()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	mux.HandleFunc("GET /arrivals/{id}", handler.ArrivalsHandler)
	mux.HandleFunc("GET /busesinservice", handler.BusesinserviceHandler)
	mux.HandleFunc("GET /vehicleinfo/{id}", handler.VehicleinfoHandler)
	mux.HandleFunc("GET /linelist", handler.LinelistHandler)
	mux.HandleFunc("GET /modelslist", handler.ModelslistHandler)
	mux.HandleFunc("GET /stoplist", handler.StoplistHandler)
	mux.HandleFunc("GET /routecodes", handler.RoutecodesHandler)
	//TODO /routestops/{id}
	mux.HandleFunc("GET /nextstops/{id}", handler.NextstopsHandler)
	mux.HandleFunc("GET /allnews", handler.AllnewsHandler)
	mux.HandleFunc("GET /news", handler.NewsHandler)
	mux.HandleFunc("GET /lineproblems", handler.LineproblemsHandler)
	mux.HandleFunc("GET /lineproblems/{id}", handler.LineproblemHandler)
	mux.HandleFunc("GET /timetable", handler.TimetableHandler)

	//Listen on port and start API
	fmt.Println("Server started on port " + config.PORT)
	log.Print(http.ListenAndServe(config.PORT, mux))
}
