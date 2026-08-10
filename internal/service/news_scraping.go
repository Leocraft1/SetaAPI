package service

import (
	"io"
	"net/http"
	"setaapi/internal/model/news"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetAllNews(url string) (news.AllNewsResponse, error) {
	response, err := http.Get(url)
	if err != nil {
		return news.AllNewsResponse{}, err
	}
	defer response.Body.Close()

	return scrapeAllNews(response.Body)
}

func scrapeAllNews(r io.Reader) (news.AllNewsResponse, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return news.AllNewsResponse{}, err
	}

	result := news.AllNewsResponse{
		News: []news.AllNews{},
	}

	doc.Find(".news li div div a").Each(func(i int, s *goquery.Selection) {
		item := news.AllNews{
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

func GetNews(url string) (news.News, error) {
	response, err := http.Get(url)
	if err != nil {
		return news.News{}, err
	}
	defer response.Body.Close()

	return scrapeNews(response.Body)
}

func scrapeNews(r io.Reader) (news.News, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return news.News{}, err
	}

	content, err := doc.Find(".descrizione").Html()
	if err != nil {
		return news.News{}, err
	}

	content = strings.TrimSpace(content)

	result := news.News{
		Title:   strings.TrimSpace(doc.Find(".container-title").Text()),
		Date:    strings.TrimSpace(doc.Find(".container-date-title").Text()),
		Content: content,
	}

	return result, nil
}
