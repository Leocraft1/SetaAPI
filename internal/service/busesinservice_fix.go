package service

import (
	"regexp"
	"setaapi/internal/model/busesinservice"
	"setaapi/internal/repository"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TODO: add news support when implemented
func FixBusesinservice(raw busesinservice.BusesRaw) busesinservice.Buses {
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
func parseBuses(from busesinservice.BusesRaw) busesinservice.Buses {
	var buses []busesinservice.BusRaw = from.Properties
	var parsed []busesinservice.Bus
	var out busesinservice.Buses

	for _, val := range buses {
		//Parses using ToDomain (model package)
		parsed = append(parsed, val.ToDomain())
	}

	out.Buses = parsed
	return out
}

func sortBusesByLine(buses []busesinservice.Bus) []busesinservice.Bus {
	sort.SliceStable(buses, func(i, j int) bool {
		numI := extractLineNumber(buses[i].Line)
		numJ := extractLineNumber(buses[j].Line)
		if numI != numJ {
			return numI < numJ
		}
		return buses[i].Line < buses[j].Line
	})

	numeric := make([]busesinservice.Bus, 0, len(buses))
	letters := make([]busesinservice.Bus, 0)

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

func fixModelAndRamp(bus *busesinservice.Bus) {
	model, ramp := repository.GetModelRampByCode(bus.Vehicle)
	if model != nil && *model != "" {
		bus.Model = *model
	}
	if ramp != nil && *model != "" {
		bus.Ramp = *ramp
	}
}

func fixPlate(bus *busesinservice.Bus) {
	targafix := strings.TrimSpace(bus.Plate_num)
	bus.Plate_num = targafix
}
