package scrapers

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"main/internal/configs"
	"main/internal/models"
	"strconv"
	"strings"
	"time"
)

type X1337Scraper struct{}

func (s X1337Scraper) Name() string {
	return "1337x"
}

func (s X1337Scraper) Search(query string) ([]models.Torrent, error) {
	cfg, _ := configs.GetScraperConfig("1337x")
	var torrents []models.Torrent

	c := colly.NewCollector(
		colly.AllowedDomains(cfg.Domains...),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	c.OnHTML("table.table-list tbody tr", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("td.name"))
		quality := models.ExtractQuality(name)
		language := models.DetectLanguage(name)

		torrent := models.Torrent{
			Name:       name,
			URL:        cfg.BaseURL + e.ChildAttr("td.name a:nth-child(2)", "href"),
			Size:       e.ChildText("td.size"),
			UploadDate: e.ChildText("td.coll-date"),
			Quality:    quality,
			Source:     s.Name(),
			Language:   language,
		}

		seedersStr := e.ChildText("td.seeds")
		seeders, err := strconv.Atoi(seedersStr)
		if err == nil {
			torrent.Seeders = seeders
		}

		leechersStr := e.ChildText("td.leeches")
		leechers, err := strconv.Atoi(leechersStr)
		if err == nil {
			torrent.Leechers = leechers
		}

		torrents = append(torrents, torrent)
	})

	formattedQuery := strings.ReplaceAll(query, " ", "+")
	searchURL, err := configs.BuildSearchURL("1337x", formattedQuery, "+")

	if err != nil {
		return nil, err
	}

	err = c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error visiting %s: %v", searchURL, err)
	}

	c.Wait()
	return torrents, nil
}

func GetMagnetLinkX1337(url string) (string, error) {
	cfg, _ := configs.GetScraperConfig("1337x")
	var magnetLink string

	c := colly.NewCollector(
		colly.AllowedDomains(cfg.Domains...),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.OnHTML("a[href^='magnet:']", func(e *colly.HTMLElement) {
		magnetLink = e.Attr("href")
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	c.Wait()
	return magnetLink, nil
}
