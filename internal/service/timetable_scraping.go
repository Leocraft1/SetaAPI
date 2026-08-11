package service

import (
	"fmt"
	"io"
	"net/http"
	"setaapi/internal/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeTimetable(timetableURL string) (model.TimetableResponse, error) {
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

	return scrapeTimetableHTML(resp.Body)
}

func scrapeTimetableHTML(r io.Reader) (model.TimetableResponse, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.TimetableResponse{}, fmt.Errorf("parse HTML: %w", err)
	}
	var result model.TimetableResponse

	var stops []string

	doc.Find(`section.lineedyn_tabella_oraria_section`).Each(
		func(sectionIndex int, section *goquery.Selection) {
			dataPercorso, exists := section.
				Find("td[data-percorso]").
				First().
				Attr("data-percorso")

			if !exists || dataPercorso == "" {
				return
			}

			var timetable *model.Timetable

			for i := range result.Timetable {
				if result.Timetable[i].DataPercorso == dataPercorso {
					timetable = &result.Timetable[i]
					break
				}
			}

			if timetable == nil {

				stops = []string{}

				section.Find(`td[data-title="Percorso"]`).Each(
					func(_ int, cell *goquery.Selection) {
						text := strings.TrimSpace(cell.Text())

						if text == "" || strings.HasPrefix(text, "VEDI IN MAPPA") {
							return
						}

						stops = append(stops, text)
					},
				)

				result.Timetable = append(
					result.Timetable,
					model.Timetable{
						DataPercorso: dataPercorso,
						Stops:        append([]string(nil), stops...),
						Timetable:    [][]string{},
					},
				)

				timetable = &result.Timetable[len(result.Timetable)-1]
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

					text := strings.TrimSpace(cell.Text())

					cells = append(cells, text)
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

				if len(column) > 0 {
					timetable.Timetable = append(
						timetable.Timetable,
						column,
					)
				}
			}
		},
	)

	return result, nil
}
