package service

import (
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"strings"
)

func fixLineInfo(val *model.LineInfo) {
	//Separates route code from journey code
	rc := val.Journey_code
	rc_split := strings.Split(rc, "-")
	rc = rc_split[0]+"-"+rc_split[1]+"-"+rc_split[2]

	//Gets specific route from DB
	linea_disp := repository.GetRoute(rc)

	//Fixes linea only if parameters are not null in DB
	if linea_disp.Disp_linea != nil && *linea_disp.Disp_linea != "" {
		val.Line = *linea_disp.Disp_linea
	}
	if linea_disp.Disp_dest != nil && *linea_disp.Disp_dest != "" {
		val.Destination = *linea_disp.Disp_dest
	}

	//Fix line_type
	switch val.Line_type {
	case "UR":
		val.Line_type = "Urbano"
	case "EX":
		val.Line_type = "Extraurbano"
	}
}

func addOfficialLine(val *model.LineInfo) {
	//Fills official_line (used for news)
	val.Official_line = val.Line
}
