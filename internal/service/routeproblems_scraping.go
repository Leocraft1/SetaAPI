package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"setaapi/internal/model/news"
	"setaapi/internal/model/routeproblems"

	"github.com/PuerkitoBio/goquery"
)

func GetRouteProblems(url string) ([]routeproblems.Problem, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("SETA returned status %d", response.StatusCode)
	}

	return ScrapeRouteProblems(response.Body)
}

func ScrapeRouteProblems(r io.Reader) ([]routeproblems.Problem, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	problems := make([]routeproblems.Problem, 0)

	doc.Find("#tabella_elenco_linee tbody tr").Each(func(i int, s *goquery.Selection) {
		numero := strings.TrimSpace(s.Find(".numero").Text())

		haProblemi := strings.TrimSpace(s.Find(".news_linea").Text()) != ""

		codiceStr, _ := s.Find(".testo_per_ricerca").Attr("data-val")
		codice, _ := strconv.Atoi(strings.TrimSpace(codiceStr))

		problems = append(problems, routeproblems.Problem{
			Num:         numero,
			HasProblems: haProblemi,
			SiteCode:    codice,
		})
	})

	return problems, nil
}

func ScrapeRouteNews(url string) (news.AllNewsResponse, error) {
	response, err := http.Get(url)
	if err != nil {
		return news.AllNewsResponse{}, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return news.AllNewsResponse{}, err
	}

	result := news.AllNewsResponse{
		News: []news.AllNews{},
	}

	doc.Find("a.news-archive-card-link").Each(func(i int, s *goquery.Selection) {
		item := news.AllNews{
			Title: strings.TrimSpace(s.Find(".title").Text()),
			Date:  strings.TrimSpace(s.Find(".date-title").Text()),
		}

		link, exists := s.Attr("href")
		if exists {
			item.Link = strings.TrimSpace(link)
		}

		result.News = append(result.News, item)
	})

	return result, nil
}

func GetSiteCode(problems []routeproblems.Problem, route string) int {
	for _, problem := range problems {
		if problem.Num == route {
			return problem.SiteCode
		}
	}

	return 0
}
