package service

import (
	"fmt"
	"setaapi/internal/model"
	"strconv"
	"strings"
)

// TODO: add news support when implemented
func FixArrivals(raw model.ArrivalRaw) model.Arrival {
	var out = parseArrivals(raw)

	//Filters planned and realtime routes duplication (vibe-coded because couldn't figure out how)
	// Step 1: costruisci il set dei Route_code che hanno una corsa realtime
	hasRealtime := make(map[string]bool)
	for _, s := range out.Arrival.Services {
		if s.State == "realtime" {
			hasRealtime[s.Journey_code] = true
		}
	}

	// Step 2: filtra
	filtered := make([]model.Service, 0, len(out.Arrival.Services))
	for _, s := range out.Arrival.Services {
		if s.State == "realtime" {
			filtered = append(filtered, s)
		} else if !hasRealtime[s.Journey_code] {
			filtered = append(filtered, s)
		}
	}

	//Calculates delay
	// Step 1: mappa le corse "planned" per Route_code
	plannedByRouteCode := make(map[string]model.Service)
	for _, s := range out.Arrival.Services {
		if s.State == "planned" {
			plannedByRouteCode[s.Journey_code] = s
		}
	}

	out.Arrival.Services = filtered

	// Step 2: per ogni corsa realtime, cerca il planned corrispondente e calcola il delay
	for i := range out.Arrival.Services {
		val := &out.Arrival.Services[i]

		if val.State != "realtime" {
			continue
		}

		planned, exists := plannedByRouteCode[val.Journey_code]
		if !exists {
			continue
		}

		delay, err := computeDelay(planned.Arrival_time, val.Arrival_time)
		if err != nil {
			continue
		}
		val.Delay = &delay
	}

	//Fix/add variants and incorrect data
	for idx := range out.Arrival.Services {
		val := &out.Arrival.Services[idx]

		addOfficialLine(&val.LineInfo)
		fixLineInfo(&val.LineInfo)
	}

	return out
}

// ---------------------
// - PRIVATE FUNCTIONS -
// ---------------------
func parseArrivals(from model.ArrivalRaw) model.Arrival {
	var services []model.ServiceRaw = from.Arrival.Services
	var parsed []model.Service
	var out model.Arrival

	for _, val := range services {
		//Parses using ToDomain (model package)
		parsed = append(parsed, val.ToDomain())
	}

	out.Arrival.Services = parsed
	return out
}

func computeDelay(plannedTime, realtimeTime string) (int, error) {
    pMinutes, err := timeStringToMinutes(plannedTime)
    if err != nil {
        return 0, err
    }
    rMinutes, err := timeStringToMinutes(realtimeTime)
    if err != nil {
        return 0, err
    }
    return rMinutes - pMinutes, nil
}

func timeStringToMinutes(t string) (int, error) {
    parts := strings.Split(t, ":")
    if len(parts) != 2 {
        return 0, fmt.Errorf("formato orario non valido: %s", t)
    }
    hours, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, err
    }
    minutes, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, err
    }
    return hours*60 + minutes, nil
}