package service

import (
	"setaapi/internal/model/busesinservice"
	//"strings"
)

// TODO: add news support when implemented
func FixBusesinservice(raw busesinservice.BusesRaw) busesinservice.Buses {
	var out = parseBuses(raw)

	/*
	for idx := range out.Buses {
		val := &out.Buses[idx]

		addOfficialLine(val)
		fixLineInfo(val)
	}*/

	return out
}

// -----------------------
// - PROTECTED FUNCTIONS -
// -----------------------
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

/*
func fixLineInfo(val *arrivals.Service) {
	//Fix "line" parameter and destination for route variants
	//Sant'Anna (Dislessia)
	if val.Destination == "SANT  ANNA" {
		val.Destination = "SANT'ANNA"
	}
	//S.Caterina (Dislessia)
	if val.Destination == "S. CATERINA" {
		val.Destination = "S.CATERINA"
	}
	//D'Avia (Dislessia)
	if val.Destination == "D AVIA" {
		val.Destination = "D'AVIA"
	}
	//La torre (Dislessia)
	if val.Destination == "L A TORRE" {
		val.Destination = "LA TORRE"
	}
	//Ragazzi Del 99 (Dislessia)
	if val.Destination == "RAGAZZI DEL   99" {
		val.Destination = "RAGAZZI DEL 99"
	}
	//1A Modena Est
	if val.Line == "1" && val.Destination == "MODENA EST" {
		val.Line = "1A"
	}
	//1A Polo Leonardo
	if val.Line == "1" && val.Destination == "POLO LEONARDO" {
		val.Line = "1A"
	}
	//1 V.ZETA - ARIETE -> ARIETE
	if val.Line == "1" && val.Destination == "V. ZETA - ARIETE" {
		val.Destination = "ARIETE"
	}
	//1/ VILLAGGIO Z
	if val.Line == "1" && val.Destination == "VILLAGGIO Z" {
		val.Line = "1/"
	}
	//1/ Autostazione (Scuola)
	if val.Line == "1" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "1/"
	}
	//1 _ -> Marinuzzi (Scuola)
	if val.Line == "1" && val.Destination == "_" {
		val.Destination = "MARINUZZI"
	}
	//2A San Donnino
	if val.Line == "2" && val.Destination == "SAN DONNINO" {
		val.Line = "2A"
	}
	//2/ Autostazione
	if val.Line == "2" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "2/"
	}
	//3A SANTA CATERINA-MONTEFIORINO (as 25/26)
	if val.Line == "3" && strings.Contains(val.Journey_code, "339-") {
		val.Line = "3A"
		val.Destination = "S.CATERINA-MONTEFIORINO"
	}
	//3A SCUOLE MARCONI-MONTEFIORINO (Scuola)
	if val.Line == "3" && strings.Contains(val.Journey_code, "289-") {
		val.Line = "3A"
		val.Destination = "SCUOLE MARCONI-MONTEFIORINO"
	}
	//3A Vaciglio
	if val.Line == "3" && val.Destination == "VACIGLIO" {
		val.Line = "3A"
	}
	//3A SCUOLE MARCONI-VACIGLIO (Scuola)
	if val.Line == "3A" && strings.Contains(val.Journey_code, "294-") {
		val.Destination = "SCUOLE MARCONI-VACIGLIO"
	}
	//3B Ragazzi del 99 (as 25/26)
	if val.Line == "3" && val.Destination == "RAGAZZI DEL 99" {
		val.Line = "3B"
	}
	//3B SCUOLE MARCONI-RAGAZZI DEL 99 (Scuola)
	if val.Line == "3B" && strings.Contains(val.Journey_code, "296-") {
		val.Destination = "SCUOLE MARCONI-RAGAZZI DEL 99"
	}
	//3B Nonantolana 1010 (as 25/26)
	if val.Line == "3" && val.Destination == "NONANTOLANA 1010" {
		val.Line = "3B"
	}
	//3B SCUOLE MARCONI-NONANTOLANA 1010 (Scuola)
	if val.Line == "3B" && strings.Contains(val.Journey_code, "287-") {
		val.Destination = "SCUOLE MARCONI-NONANTOLANA 1010"
	}
	//3/ Stazione FS (as 25/26)
	if val.Line == "3" && val.Destination == "STAZIONE FS" {
		val.Line = "3/"
	}
	//3A MONTEFIORINO (Domenica)
	if val.Line == "3" && (strings.Contains(val.Journey_code, "407-") || strings.Contains(val.Journey_code, "327-")) {
		val.Line = "3A"
		val.Destination = "MONTEFIORINO"
	}
	//4 POLO LEONARDO-GALILEI
	if val.Line == "4" && strings.Contains(val.Journey_code, "MO4-As-432") {
		val.Destination = "POLO LEON.-GALILEI"
	}
	//4/ Autostazione (as 25/26)
	if val.Line == "4" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "4/"
	}
	//4/ STAZIONE FS
	if val.Line == "4" && val.Destination == "STAZIONE FS" {
		val.Line = "4/"
	}
	//5 Dalla Chiesa -> La Torre
	if val.Line == "5" && val.Destination == "DALLA CHIESA" {
		val.Destination = "LA TORRE"
	}
	//5A Tre Olmi
	if val.Line == "5" && val.Destination == "TRE OLMI" {
		val.Line = "5A"
	}
	//6A Santi (as 25/26)
	if val.Line == "6" && val.Destination == "SANTI" {
		val.Line = "6A"
	}
	//6B Villanova (as 25/26)
	if val.Line == "6" && val.Destination == "VILLANOVA" {
		val.Line = "6B"
	}
	//7 GOTTARDI -> POLICLINICO GOTTARDI
	if val.Line == "7" && val.Destination == "GOTTARDI" {
		val.Destination = "POLICLINICO GOTTARDI"
	}
	//7/ Stazione FS
	if val.Line == "7" && val.Destination == "STAZIONE FS" {
		val.Line = "7/"
	}
	//7/ AUTOSTAZIO§NE
	if val.Line == "7" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "7/"
	}
	//7A STAZIONE FS -> GOTTARDI
	if val.Line == "7A" && val.Destination == "STAZIONE FS" && !strings.Contains(val.Journey_code, "728-") {
		val.Destination = "GOTTARDI"
	}
	//7A/ STAZIONE FS
	if val.Line == "7A" && val.Destination == "STAZIONE FS" {
		val.Line = "7A/"
	}
	//9A Marzaglia Nuova
	if val.Line == "9" && val.Destination == "MARZAGLIA" {
		val.Line = "9A"
		val.Destination = "MARZAGLIA NUOVA"
	}
	//9A MARZAGLIA NUOVA
	if val.Line == "9" && val.Destination == "MARZAGLIA NUOVA" {
		val.Line = "9A"
	}
	//9B VIRGILIO
	if val.Line == "9" && val.Destination == "VIRGILIO" {
		val.Line = "9B"
	}
	//9B BRAGHIROLI-GOTTARDI
	if val.Line == "9" && strings.Contains(val.Journey_code, "MO9-As-948") {
		val.Line = "9B"
		val.Destination = "BRAGHIROLI-GOTTARDI"
	}
	//9C Rubiera
	if val.Line == "9" && val.Destination == "RUBIERA" {
		val.Line = "9C"
	}
	//9/ Stazione FS
	if val.Line == "9" && val.Destination == "STAZIONE FS" {
		val.Line = "9/"
	}
	//9/ Autostazione
	if val.Line == "9" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "9/"
	}
	//10A La Rocca
	if val.Line == "10" && val.Destination == "LA ROCCA" {
		val.Line = "10A"
	}
	//10B MARZAGLIA NUOVA
	if val.Line == "10" && val.Destination == "MARZAGLIA NUOVA" {
		val.Line = "10B"
		val.Destination = "COGNENTO-MARZAGLIA NUOVA"
	}
	//10/ Liceo Sigonio
	if val.Line == "10" && val.Destination == "LICEO SIGONIO" {
		val.Line = "10/"
	}
	//10/ POLO LEONARDO
	if val.Line == "10" && val.Destination == "POLO LEONARDO" {
		val.Line = "10/"
	}
	//10/ AUTOSTAZIONE
	if val.Line == "10" && val.Destination == "AUTOSTAZIONE" {
		val.Line = "10/"
	}
	//11/ Stazione FS
	if val.Line == "11" && val.Destination == "STAZIONE FS" {
		val.Line = "11/"
	}
	//12A Nazioni ma sono dei coglioni di merda
	if val.Line == "12" && (strings.Contains(val.Journey_code, "1280-") || strings.Contains(val.Journey_code, "1284-")) {
		val.Line = "12A"
		val.Destination = "NAZIONI"
	}
	//12/ Fanti FS
	if val.Line == "12" && val.Destination == "FANTI FS" {
		val.Line = "12/"
	}
	//12/ Garibaldi (Scuola)
	if val.Line == "12" && val.Destination == "GARIBALDI" {
		val.Line = "12/"
	}
	//12/ Largo Garibaldi (Scuola)
	if val.Line == "12" && val.Destination == "LARGO GARIBALDI" {
		val.Line = "12/"
	}
	//13 S.ANNA -> SANT'ANNA (dio rincoglionito e dislessico)
	if val.Line == "13" && val.Destination == "S.ANNA" {
		val.Destination = "SANT'ANNA"
	}
	//13A Carcere
	if val.Line == "13" && val.Destination == "CARCERI" {
		val.Line = "13A"
	}
	//13F Variante di merda
	if val.Line == "13" && strings.Contains(val.Journey_code, "133") {
		val.Line = "13F"
	}
	//394 VIA GIANNONE -> CINEMA ESTIVO
	if val.Line == "394" && val.Destination == "VIA GIANNONE" {
		val.Destination = "CINEMA ESTIVO"
	}
	//643 _ -> Polo Scolastico Sassuolo
	if val.Line == "643" && val.Destination == "_" {
		val.Destination = "POLO SCOLASTICO SASSUOLO"
	}
}

func addOfficialLine(val *arrivals.Service) {
	//Fills official_line (used for news)
	val.Official_line = val.Line
}
*/