package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"setaapi/config"
	"setaapi/internal/model/arrivals"
	"setaapi/internal/model/busesinservice"
	"setaapi/internal/model/routeproblems"
	"setaapi/internal/service"
	"strconv"
)

// URLs decl. section
var arrivalsBaseUrl string = "https://avm.setaweb.it/SETA_WS/services/arrival/"
var wimbBaseUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO"
var nextstopsUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/getwaypointarrivals/"
var newsUrl string = "https://www.setaweb.it/mo/news"
var lineeDynUrl string = "https://www.setaweb.it/mo/lineedyn"
var lineeNewsUrl string = "https://www.setaweb.it/mo/news/linea/"

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
	response, err := http.Get("http://localhost" + config.PORT + "/busesinservice")
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

// GET /linelist
func LinelistHandler(w http.ResponseWriter, r *http.Request) {
	nums := service.GetRoutenums()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(nums)
}

// GET /modelslist
func ModelslistHandler(w http.ResponseWriter, r *http.Request) {
	nums := service.GetModels()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(nums)
}

// GET /stoplist
func StoplistHandler(w http.ResponseWriter, r *http.Request) {
	stops := service.GetStops()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(stops)
}

//TODO GET /routestops/{id}

// GET /nextstops/{id}
func NextstopsHandler(w http.ResponseWriter, r *http.Request) {
	journey_code := r.PathValue("id")
	response, err := http.Get(nextstopsUrl + journey_code)
	if err != nil {
		fmt.Println("Error fetching next stops", err)
	}
	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	//Parses response back to JSON and returns it to client
	w.Write(body)
}

// GET /allnews
// Collects all the news scraped from SETA's website
func AllnewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := http.Get(newsUrl)
	if err != nil {
		fmt.Println("AllnewsHandler error: ", err)
		w.Write([]byte("AllnewsHandler error: " + err.Error()))
	}

	defer response.Body.Close()

	w.WriteHeader(response.StatusCode)

	service.Scrapeallnews(response.Body, w)
}

// GET /news
// Scrapes the news page given by the URI link parameter
func NewsHandler(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Query().Get("link")

	if link == "" {
		http.Error(w, "missing link parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := http.Get(link)
	if err != nil {
		fmt.Println("NewsHandler error: ", err)
		w.Write([]byte("NewsHandler error: " + err.Error()))
	}

	defer response.Body.Close()

	w.WriteHeader(response.StatusCode)

	service.Scrapenews(response.Body, w)
}

// GET /routeproblems
// Collects all the route problems scraped from SETA's website page
func RouteproblemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	problems, err := service.GetRouteProblems(lineeDynUrl)
	if err != nil {
		http.Error(w, "RouteproblemsHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	response := routeproblems.ProblemCodesResponse{
		Problem: problems,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("RouteproblemsHandler JSON error:", err)
	}
}

// GET /routeproblems/{id}
// Collects the news from the given route scraped from SETA's website page
func RouteproblemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	route := r.PathValue("id")

	problems, err := service.GetRouteProblems(lineeDynUrl)
	if err != nil {
		http.Error(w, "RouteproblemHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	siteCode := service.GetSiteCode(problems, route)

	news, err := service.ScrapeRouteNews(lineeNewsUrl + strconv.Itoa(siteCode))
	if err != nil {
		http.Error(w, "RouteproblemHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(news); err != nil {
		fmt.Println("RouteproblemHandler JSON error:", err)
	}
}
