package service

import (
	"fmt"
	"io"
	"net/http"
	"setaapi/internal/model"
	"setaapi/internal/repository"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeTimetable(timetableURL string, linea string, verso string) (model.TimetableResponse, error) {
	resp, err := http.Get(timetableURL)
	if err != nil {
		return model.TimetableResponse{}, fmt.Errorf("fetch timetable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.TimetableResponse{}, fmt.Errorf(
			"SETA returned status %d",
			resp.StatusCode,
		)
	}

	scraped, err := scrapeTimetableHTML(resp.Body, linea, verso)
	fixed := fixTimetable(scraped)
	return fixed, err
}

func scrapeTimetableHTML(r io.Reader, linea string, verso string) (model.TimetableResponse, error) {
    doc, err := goquery.NewDocumentFromReader(r)
    if err != nil {
        return model.TimetableResponse{}, fmt.Errorf("parse HTML: %w", err)
    }

    var result model.TimetableResponse
    stopsCaptured := false

    doc.Find(`section.lineedyn_tabella_oraria_section`).Each(func(_ int, section *goquery.Selection) {
        // Un data-percorso per OGNI colonna dell'header, non uno solo per sezione
        var routeCodes []string
        section.Find(`thead td[data-percorso]`).Each(func(_ int, cell *goquery.Selection) {
            code, _ := cell.Attr("data-percorso")
            routeCodes = append(routeCodes, code)
        })
        if len(routeCodes) == 0 {
            return
        }

        if !stopsCaptured {
            section.Find(`td[data-title="Percorso"]`).Each(func(_ int, cell *goquery.Selection) {
                text := strings.TrimSpace(cell.Text())
                if text == "" || strings.HasPrefix(text, "VEDI IN MAPPA") {
                    return
                }
                result.Stops = append(result.Stops, text)
            })
            stopsCaptured = true
        }

        rows := section.Find("tbody tr")
        if rows.Length() == 0 {
            return
        }

        var rowsData [][]string
        rows.Each(func(_ int, row *goquery.Selection) {
            var cells []string
            row.Find("td").Each(func(_ int, cell *goquery.Selection) {
                if cell.Is(`[data-title="Percorso"]`) {
                    return
                }
                cells = append(cells, strings.TrimSpace(cell.Text()))
            })
            if len(cells) > 0 {
                rowsData = append(rowsData, cells)
            }
        })
        if len(rowsData) == 0 {
            return
        }

        columnCount := len(rowsData[0])
        for columnIndex := 0; columnIndex < columnCount; columnIndex++ {
            column := make([]string, 0, len(rowsData))
            for _, row := range rowsData {
                if columnIndex >= len(row) {
                    continue
                }
                column = append(column, row[columnIndex])
            }
            if len(column) == 0 {
                continue
            }

            routeCode := ""
            if columnIndex < len(routeCodes) {
                routeCode = "MO"+ linea + "-" + verso + "-" +routeCodes[columnIndex] // codice specifico DI QUESTA colonna
            }

            result.Journeys = append(result.Journeys, model.Journey{
                RouteCode: routeCode,
                Times:     column,
            })
        }
    })

    return result, nil
}

//Adds display specs to response
func fixTimetable(timetable model.TimetableResponse) model.TimetableResponse {
	rcs := repository.GetRCTable()
	//Search for matching route_code
	for idx1 := range timetable.Journeys {
		val1 := &timetable.Journeys[idx1]
		for _, val2 := range rcs {
			if val1.RouteCode == val2.Rc {
				val1.Disp_line = val2.Disp_linea
				val1.Disp_dest = val2.Disp_dest
				break
			}
		}
	}

	return timetable
}