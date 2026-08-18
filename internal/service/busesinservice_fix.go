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

func GetBusesInservice(url string, problemsUrl string) (model.Buses, error) {
	response, err := http.Get(url)
	if err != nil {
		return model.Buses{}, err
	}
	defer response.Body.Close()

	// News section
	problemsRaw, err := http.Get(problemsUrl)
	if err != nil {
		return model.Buses{}, err
	}
	defer problemsRaw.Body.Close()

	problems, err := ScrapeRouteProblems(problemsRaw.Body)
	if err != nil {
		return model.Buses{}, err
	}

	// AEP support
	aepRaw := repository.GetAEP()

	var aepMap = make(map[int]bool)
	for _, val := range aepRaw {
		aepMap[val.Matricola] = val.Has_aep
	}

	var raw model.BusesRaw
	if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
		return model.Buses{}, err
	}

	return FixBusesinservice(raw, model.ProblemCodesResponse{
		Problems: problems,
	}, aepMap), nil
}

func FixBusesinservice(raw model.BusesRaw, problems model.ProblemCodesResponse, aep map[int]bool) model.Buses {
	var out = parseBuses(raw)

	//Sorts buses
	out.Buses = sortBusesByLine(out.Buses)
	for idx := range out.Buses {
		val := &out.Buses[idx]

		addOfficialLine(&val.LineInfo)
		fixLineInfo(&val.LineInfo)
		fixModelAndRamp(val)
		fixPlate(val)
		vehicleInt, _ := strconv.ParseInt(val.Vehicle, 0, 64)
		val.Has_AEP = aep[int(vehicleInt)]
	}

	//News section
	for idx1 := range out.Buses {
		val1 := &out.Buses[idx1]
		for _, val2 := range problems.Problems {
			if val1.Official_line == val2.Num && val2.HasProblems {
				val1.Has_problems = true
				break
			} else if val1.Official_line == val2.Num && !val2.HasProblems {
				//If there are no problems break the cycle
				break
			}
		}
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
