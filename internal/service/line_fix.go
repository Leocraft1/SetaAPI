package service

import (
	"setaapi/internal/model"
	"setaapi/internal/repository"
)

func fixLineInfo(val *model.LineInfo) {
	//Separates route code from journey code
	rc := val.Journey_code
	rc = rc[:len(rc)-8]

	//Gets specific route from DB
	linea_disp := repository.GetRoute(rc)
	
	//Fixes linea only if parameters are not null in DB
	if linea_disp.Disp_linea != nil {
		val.Line = *linea_disp.Disp_linea
	}
	if linea_disp.Disp_dest != nil {
		val.Destination = *linea_disp.Disp_dest
	}
}

func addOfficialLine(val *model.LineInfo) {
	//Fills official_line (used for news)
	val.Official_line = val.Line
}
