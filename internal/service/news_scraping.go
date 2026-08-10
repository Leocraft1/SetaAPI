package service

import (
	"io"
	"net/http"
	"setaapi/internal/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeAllNews(r io.Reader) (model.AllNewsResponse, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.AllNewsResponse{}, err
	}

	result := model.AllNewsResponse{
		News: []model.AllNews{},
	}

	doc.Find(".news li div div a").Each(func(i int, s *goquery.Selection) {
		item := model.AllNews{
			Title: strings.TrimSpace(s.Find(".title").Text()),
			Date:  strings.TrimSpace(s.Find(".date-title").Text()),
		}

		if link, ok := s.Attr("href"); ok {
			item.Link = strings.TrimSpace(link)
		}

		if image, ok := s.Find(".image-news").Attr("data-bg"); ok {
			switch {
			case strings.Contains(image, "pericolo.png"):
				item.Type = "Importante"

			case strings.Contains(image, "seta-informa.png"):
				item.Type = "Informazione"

			case strings.Contains(image, "novita.png"):
				item.Type = "Novità"

			case strings.Contains(image, "orari.png"):
				item.Type = "Orari"

			case strings.Contains(image, "autobus-treno.png"):
				item.Type = "Autobus Treno"

			case strings.Contains(image, "tessera.png"):
				item.Type = "Biglietti"

			case strings.Contains(image, "lavori-in-corso.png"):
				item.Type = "Lavori in corso"

			case strings.Contains(image, "controllore.png"):
				item.Type = "Personale"

			case strings.Contains(image, "salvadanaio.png"):
				item.Type = "Agevolazioni"

			case strings.Contains(image, "abbonamento.png"):
				item.Type = "Abbonamenti"

			case strings.Contains(image, "scuola.png"):
				item.Type = "Scuola"
			}
		}

		result.News = append(result.News, item)
	})

	return result, nil
}

func GetNews(url string) (model.News, error) {
	response, err := http.Get(url)
	if err != nil {
		return model.News{}, err
	}
	defer response.Body.Close()

	return scrapeNews(response.Body)
}

func scrapeNews(r io.Reader) (model.News, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.News{}, err
	}

	content, err := doc.Find(".descrizione").Html()
	if err != nil {
		return model.News{}, err
	}

	content = strings.TrimSpace(content)

	result := model.News{
		Title:   strings.TrimSpace(doc.Find(".container-title").Text()),
		Date:    strings.TrimSpace(doc.Find(".container-date-title").Text()),
		Content: content,
	}

	return result, nil
}
