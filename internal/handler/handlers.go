package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"setaapi/config"
	"setaapi/internal/model"
	"setaapi/internal/service"
	"strconv"
	"time"
)

// URLs decl. section
var arrivalsBaseUrl string = "https://avm.setaweb.it/SETA_WS/services/arrival/"
var wimbBaseUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO"
var nextstopsUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/getwaypointarrivals/"
var newsUrl string = "https://www.setaweb.it/mo/news"
var lineeDynUrl string = "https://www.setaweb.it/mo/lineedyn"
var lineeNewsUrl string = "https://www.setaweb.it/mo/news/linea/"
var timetablesUrl string = "https://www.setaweb.it/mo/lineedyn/corse-tabella"

// GET /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	message := "API is healthy and running."

	w.Write([]byte(message))
}

// GET /arrivals/{id}
func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	stopId := r.PathValue("id")
	response, err := http.Get(arrivalsBaseUrl + stopId)

	//Checks for connection errors
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
	}

	//Parses response into struct
	var arrivalsRaw model.ArrivalRaw
	json.NewDecoder(response.Body).Decode(&arrivalsRaw)

	//News section
	problemsRaw, err := http.Get("http://localhost"+config.PORT+"/lineproblems")
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
	}

	var problems model.ProblemCodesResponse
	json.NewDecoder(problemsRaw.Body).Decode(&problems)

	//Fixes incorrect stuff/add parameters
	arrivals := service.FixArrivals(arrivalsRaw, problems)

	//Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(arrivals)
}

// GET /busesinservice
func BusesinserviceHandler(w http.ResponseWriter, r *http.Request) {
	buses, err := service.GetBusesInservice(wimbBaseUrl)
	if err != nil {
		http.Error(w, "BusesinserviceHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(buses)
}

// GET /vehicleinfo/{id}
// Searches for requested vehicle on /busesinservice response
func VehicleinfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	buses, err := service.GetBusesInservice(wimbBaseUrl)
	if err != nil {
		http.Error(w, "VehicleinfoHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for _, bus := range buses.Buses {
		if bus.Vehicle == id {
			json.NewEncoder(w).Encode(bus)
			return
		}
	}

	http.Error(w, "Vehicle not found.", http.StatusNotFound)
}

// GET /linelist
func LinelistHandler(w http.ResponseWriter, r *http.Request) {
	nums := service.GetRoutenums()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(nums)
}

// GET /modelslist
func ModelslistHandler(w http.ResponseWriter, r *http.Request) {
	nums := service.GetModels()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(nums)
}

// GET /stoplist
func StoplistHandler(w http.ResponseWriter, r *http.Request) {
	stops := service.GetStops()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stops)
}

// GET /routecodes
func RoutecodesHandler(w http.ResponseWriter, r *http.Request) {
	routelist := service.GetRoutecodes()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(routelist)
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

	w.Write(body)
}

// GET /allnews
func AllnewsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(newsUrl)
	if err != nil {
		http.Error(w, "AllnewsHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	result, err := service.ScrapeAllNews(response.Body)
	if err != nil {
		http.Error(w, "AllnewsHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

// GET /news?link=
func NewsHandler(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Query().Get("link")

	//Verifies if link parameter exists
	if link == "" {
		http.Error(w, "Missing link parameter", http.StatusBadRequest)
		return
	}

	result, err := service.GetNews(link)
	if err != nil {
		http.Error(w, "NewsHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

// GET /lineproblems
func LineproblemsHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := http.Get(lineeDynUrl)
	if err != nil {
		http.Error(w, "RouteproblemsHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}
	problems, err := service.ScrapeRouteProblems(raw.Body)
	if err != nil {
		http.Error(w, "RouteproblemsHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.ProblemCodesResponse{
		Problems: problems,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// GET /lineproblems/{id}
func LineproblemHandler(w http.ResponseWriter, r *http.Request) {
	route := r.PathValue("id")

	raw, err := http.Get(lineeDynUrl)
	if err != nil {
		http.Error(w, "RouteproblemsHandler error: "+err.Error(), http.StatusBadGateway)
		return
	}

	problems, err := service.ScrapeRouteProblems(raw.Body)
	if err != nil {
		http.Error(w, "RouteproblemsHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	siteCode := service.GetSiteCode(problems, route)

	news, err := service.ScrapeRouteNews(lineeNewsUrl + strconv.Itoa(siteCode))
	if err != nil {
		http.Error(w, "RouteproblemHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(news)
}

// GET /timetable
// Scrapes the line timetable from SETA's website
func TimetableHandler(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("line")
	verse := r.URL.Query().Get("verse")

	if line == "" {
		http.Error(w, "missing line parameter", http.StatusBadRequest)
		return
	}

	if verse == "" {
		http.Error(w, "missing verse parameter", http.StatusBadRequest)
		return
	}

	params := url.Values{}
	params.Set("l", "MO"+line)
	params.Set("v", verse)
	params.Set("g", time.Now().Format("2006-01-02"))

	timetableURL := timetablesUrl + "?" + params.Encode()

	response, err := service.ScrapeTimetable(timetableURL, line, verse)
	if err != nil {
		http.Error(
			w,
			"TimetableHandler error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
