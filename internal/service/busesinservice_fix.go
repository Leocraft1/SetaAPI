package service

import (
	"encoding/json"
	"net/http"
	"regexp"
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func GetBusesInservice(url string) (model.Buses, error) {
	response, err := http.Get(url)
	if err != nil {
		return model.Buses{}, err
	}
	defer response.Body.Close()

	var raw model.BusesRaw

	if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
		return model.Buses{}, err
	}

	return FixBusesinservice(raw), nil
}

// TODO: add news support when implemented
func FixBusesinservice(raw model.BusesRaw) model.Buses {
	var out = parseBuses(raw)

	//Sorts buses
	out.Buses = sortBusesByLine(out.Buses)
	for idx := range out.Buses {
		val := &out.Buses[idx]

		addOfficialLine(&val.LineInfo)
		fixLineInfo(&val.LineInfo)
		fixModelAndRamp(val)
		fixPlate(val)
	}

	return out
}

// ---------------------
// - PRIVATE FUNCTIONS -
// ---------------------
func parseBuses(from model.BusesRaw) model.Buses {
	var buses []model.BusRaw = from.Properties
	var parsed []model.Bus
	var out model.Buses

	for _, val := range buses {
		//Parses using ToDomain (model package)
		parsed = append(parsed, val.ToDomain())
	}

	out.Buses = parsed
	return out
}

func sortBusesByLine(buses []model.Bus) []model.Bus {
	sort.SliceStable(buses, func(i, j int) bool {
		numI := extractLineNumber(buses[i].Line)
		numJ := extractLineNumber(buses[j].Line)
		if numI != numJ {
			return numI < numJ
		}
		return buses[i].Line < buses[j].Line
	})

	numeric := make([]model.Bus, 0, len(buses))
	letters := make([]model.Bus, 0)

	for _, b := range buses {
		if len(b.Line) > 0 && unicode.IsLetter(rune(b.Line[0])) {
			letters = append(letters, b)
		} else {
			numeric = append(numeric, b)
		}
	}

	return append(numeric, letters...)
}

var numericPartRegex = regexp.MustCompile(`\d+`)

func extractLineNumber(line string) int {
	match := numericPartRegex.FindString(line)
	if match == "" {
		return 0
	}
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return num
}

func fixModelAndRamp(bus *model.Bus) {
	model, ramp := repository.GetModelRampByCode(bus.Vehicle)
	if model != nil && *model != "" {
		bus.Model = *model
	}
	if ramp != nil && *model != "" {
		bus.Ramp = *ramp
	}
}

func fixPlate(bus *model.Bus) {
	targafix := strings.TrimSpace(bus.Plate_num)
	bus.Plate_num = targafix
}
