package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"setaapi/internal/service"
	"strconv"
	"strings"
	"time"
)

// URLs decl. section
const (
	ArrivalsBaseUrl string = "https://avm.setaweb.it/SETA_WS/services/arrival/"
	WimbBaseUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/map/MO"
	NextstopsUrl string = "https://wimb.setaweb.it/publicmapbe/vehicles/getwaypointarrivals/"
	NewsUrl string = "https://www.setaweb.it/mo/news"
	LineeDynUrl string = "https://www.setaweb.it/mo/lineedyn"
	LineeNewsUrl string = "https://www.setaweb.it/mo/news/linea/"
	TimetablesUrl string = "https://www.setaweb.it/mo/lineedyn/corse-tabella"
	PercorsoAutistaUrl string = "https://www.setaweb.it/percorsoAutista/percorso_mappa.php"
)

func addCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// GET /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	message := "API is healthy and running."

	w.Write([]byte(message))
}

// GET /arrivals/{id}
func ArrivalsHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	stopId := r.PathValue("id")
	response, err := http.Get(ArrivalsBaseUrl + stopId)

	//Checks for connection errors
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
	}
	defer response.Body.Close()

	//Parses response into struct
	var arrivalsRaw model.ArrivalRaw
	json.NewDecoder(response.Body).Decode(&arrivalsRaw)

	//News section
	rawProblems, err := http.Get(LineeDynUrl)
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
	}
	defer response.Body.Close()

	problems, err := service.ScrapeRouteProblems(rawProblems.Body)
	if err != nil {
		fmt.Println("ArrivalsHandler error: ", err)
		w.Write([]byte("ArrivalsHandler error: " + err.Error()))
		return
	}

	// Fixes incorrect stuff/add parameters
	arrivals := service.FixArrivals(arrivalsRaw, model.ProblemCodesResponse{
		Problems: problems,
	})

	//Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	//Parses response back to JSON and returns it to client
	json.NewEncoder(w).Encode(arrivals)
}

// GET /busesinservice
func BusesinserviceHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	buses, err := service.GetBusesinservice(WimbBaseUrl, LineeDynUrl)
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
	addCORS(w)
	id := r.PathValue("id")

	buses, err := service.GetBusesinservice(WimbBaseUrl, LineeDynUrl)
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
	addCORS(w)
	nums := service.GetRouteByRCnums()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(nums)
}

// GET /modelslist
func ModelslistHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	nums := service.GetModels()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(nums)
}

// GET /stops
func StoplistHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	stops := service.GetStops()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stops)
}

// GET /stopsinfo
func StopcountHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	count, timestamp := repository.GetStopCount()

	response := model.StopsInfo{
		Count: count,
		Updated_at: timestamp,
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// GET /routecodes
func RoutecodesHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	routelist := service.GetRouteByRCcodes()

	//Set headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(routelist)
}

// GET /nextstops/{id}
func NextstopsHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	journey_code := r.PathValue("id")
	response, err := http.Get(NextstopsUrl + journey_code)
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
	addCORS(w)
	response, err := http.Get(NewsUrl)
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
	addCORS(w)
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
	addCORS(w)
	raw, err := http.Get(LineeDynUrl)
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
	addCORS(w)
	route := r.PathValue("id")

	raw, err := http.Get(LineeDynUrl)
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

	news, err := service.ScrapeRouteNews(LineeNewsUrl + strconv.Itoa(siteCode))
	if err != nil {
		http.Error(w, "RouteproblemHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(news)
}

// GET /timetable
func TimetableHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	line := r.URL.Query().Get("line")
	verse := r.URL.Query().Get("verse")

	if line == "" {
		http.Error(w, "Missing line parameter", http.StatusBadRequest)
		return
	}

	if verse == "" {
		http.Error(w, "Missing verse parameter", http.StatusBadRequest)
		return
	}

	params := url.Values{}
	params.Set("l", "MO"+line)
	params.Set("v", verse)
	params.Set("g", time.Now().Format("2006-01-02"))

	timetableURL := TimetablesUrl + "?" + params.Encode()

	response, err := service.ScrapeTimetable(timetableURL, line, verse)
	if err != nil {
		http.Error(w, "TimetableHandler error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// GET /routestops/{id}
func RoutestopsHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	route_code := r.PathValue("id")

	out := service.GetRoutestops(route_code)

	//Sets headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(out)
}

// GET /routemap/{id}
func RoutemapHandler(w http.ResponseWriter, r *http.Request) {
	addCORS(w)
	code := r.PathValue("id")

	now := time.Now()
	todayDate := fmt.Sprintf("%d-%d-%d", now.Year(), int(now.Month()), now.Day())

	formData := url.Values{}
	formData.Set("data", todayDate)
	formData.Set("percorso", code)

	req, err := http.NewRequest("POST", PercorsoAutistaUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fixedURLs, err := service.FixRelativeUrls(string(bodyBytes), "https://www.setaweb.it/percorsoAutista/")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	finalHTML := service.FixLeafletScript(service.HideTopBar(fixedURLs))

	//Set headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(finalHTML))
}
