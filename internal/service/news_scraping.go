package service

import (
	"encoding/json"
	"fmt"
	"io"
	"setaapi/internal/model/news"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrapeallnews(r io.ReadCloser, w io.Writer) {
	defer r.Close()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Scrapeallnews error:", err)
		return
	}

	result := news.AllNewsResponse{
		News: []news.AllNews{},
	}

	doc.Find(".news li div div a").Each(func(i int, s *goquery.Selection) {
		news := news.AllNews{
			Title: strings.TrimSpace(s.Find(".title").Text()),
			Date:  strings.TrimSpace(s.Find(".date-title").Text()),
		}

		if link, ok := s.Attr("href"); ok {
			news.Link = strings.TrimSpace(link)
		}

		if image, ok := s.Find(".image-news").Attr("data-bg"); ok {
			switch {
			case strings.Contains(image, "pericolo.png"):
				news.Type = "Importante"

			case strings.Contains(image, "seta-informa.png"):
				news.Type = "Informazione"

			case strings.Contains(image, "novita.png"):
				news.Type = "Novità"

			case strings.Contains(image, "orari.png"):
				news.Type = "Orari"

			case strings.Contains(image, "autobus-treno.png"):
				news.Type = "Autobus Treno"

			case strings.Contains(image, "tessera.png"):
				news.Type = "Biglietti"

			case strings.Contains(image, "lavori-in-corso.png"):
				news.Type = "Lavori in corso"

			case strings.Contains(image, "controllore.png"):
				news.Type = "Personale"

			case strings.Contains(image, "salvadanaio.png"):
				news.Type = "Agevolazioni"

			case strings.Contains(image, "abbonamento.png"):
				news.Type = "Abbonamenti"

			case strings.Contains(image, "scuola.png"):
				news.Type = "Scuola"
			}
		}

		result.News = append(result.News, news)
	})

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(result); err != nil {
		fmt.Println("Scrapeallnews JSON error:", err)
	}
}

func Scrapenews(r io.ReadCloser, w io.Writer) {
	defer r.Close()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Scrapenews error:", err)
		return
	}

	content, _ := doc.Find(".descrizione").Html()
	content = strings.TrimSpace(content)

	result := news.News{
		Title:   strings.TrimSpace(doc.Find(".container-title").Text()),
		Date:    strings.TrimSpace(doc.Find(".container-date-title").Text()),
		Content: content,
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(result); err != nil {
		fmt.Println("Scrapenews JSON error:", err)
	}
}
