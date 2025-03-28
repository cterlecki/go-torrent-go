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

type RARBGScraper struct{}

func (s RARBGScraper) Name() string {
	return "RARBG"
}

func (s RARBGScraper) Search(query string) ([]models.Torrent, error) {
	cfg, _ := configs.GetScraperConfig("rarbg")
	var torrents []models.Torrent

	c := colly.NewCollector(
		colly.AllowedDomains(cfg.Domains...),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	c.OnHTML("table.lista2t tr.lista2", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("td:nth-child(2) a"))
		if name == "" {
			return
		}

		quality := models.ExtractQuality(name)
		language := models.DetectLanguage(name)

		torrent := models.Torrent{
			Name:     name,
			URL:      cfg.BaseURL + e.ChildAttr("td:nth-child(2) a", "href"),
			Size:     e.ChildText("td:nth-child(4)"),
			Quality:  quality,
			Source:   s.Name(),
			Language: language,
		}

		seedersStr := e.ChildText("td:nth-child(5)")
		seeders, err := strconv.Atoi(seedersStr)
		if err == nil {
			torrent.Seeders = seeders
		}

		leechersStr := e.ChildText("td:nth-child(6)")
		leechers, err := strconv.Atoi(leechersStr)
		if err == nil {
			torrent.Leechers = leechers
		}

		torrents = append(torrents, torrent)
	})

	formattedQuery := strings.ReplaceAll(query, " ", "+")
	searchURL, err := configs.BuildSearchURL("rarbg", formattedQuery, "+")

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
