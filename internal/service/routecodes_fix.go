package service

import (
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"sort"
	"unicode"
)

func GetRoutecodes() model.RLResponse {
	rcs := repository.GetRCTable()

	type lineGroup struct {
		active []model.RCEntry // Still_exists == 0
		old    []model.RCEntry // Still_exists == 1
	}

	grouped := make(map[string]*lineGroup)
	var lineOrder []string

	for _, val := range rcs {
		g, exists := grouped[val.Linea]
		if !exists {
			g = &lineGroup{}
			grouped[val.Linea] = g
			lineOrder = append(lineOrder, val.Linea)
		}

		entry := model.RCEntry{
			//Result parameters
			Rc:        val.Rc,
			Desc:      val.Desc,
			Disp_line: val.Disp_linea,
			Disp_dest: val.Disp_dest,
			Exists:    val.Still_exists,
		}

		if val.Still_exists == false {
			g.old = append(g.old, entry)
		} else {
			g.active = append(g.active, entry)
		}
	}

	result := make([]model.RouteListElement, 0, len(lineOrder))
	for _, line := range lineOrder {
		g := grouped[line]
		result = append(result, model.RouteListElement{
			Line:        line,
			Route_codes: append(g.active, g.old...),
		})
	}

	//Sort by extracting line number
	result = sortRC(result)

	response := model.RLResponse{List: result}

	return response
}

func sortRC(buses []model.RouteListElement) []model.RouteListElement {
	sort.SliceStable(buses, func(i, j int) bool {
		numI := extractLineNumber(buses[i].Line)
		numJ := extractLineNumber(buses[j].Line)
		if numI != numJ {
			return numI < numJ
		}
		return buses[i].Line < buses[j].Line
	})

	numeric := make([]model.RouteListElement, 0, len(buses))
	letters := make([]model.RouteListElement, 0)

	for _, b := range buses {
		if len(b.Line) > 0 && unicode.IsLetter(rune(b.Line[0])) {
			letters = append(letters, b)
		} else {
			numeric = append(numeric, b)
		}
	}

	return append(numeric, letters...)
}
