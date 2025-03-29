package scrapers

import (
	"fmt"
	"github.com/cterlecki/go-torrent-go/internal/configs"
	"github.com/cterlecki/go-torrent-go/internal/models"
	"github.com/gocolly/colly/v2"
	"strings"
	"time"
)

type Torrent9Scraper struct{}

func (s Torrent9Scraper) Name() string {
	return "Torrent9"
}

func (s Torrent9Scraper) Search(query string) ([]models.Torrent, error) {
	cfg, _ := configs.GetScraperConfig("torrent9")
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

	c.OnHTML("table.table tbody tr", func(e *colly.HTMLElement) {
		nameLink := e.ChildAttr("td a", "href")
		if nameLink == "" {
			return
		}

		name := strings.TrimSpace(e.ChildText("td a"))
		quality := models.ExtractQuality(name)
		language := models.DetectLanguage(name)

		torrent := models.Torrent{
			Name:     name,
			URL:      cfg.BaseURL + nameLink,
			Size:     e.ChildText("td:nth-child(2)"),
			Seeders:  parseInt(e.ChildText("td:nth-child(3)")),
			Leechers: parseInt(e.ChildText("td:nth-child(4)")),
			Quality:  quality,
			Source:   s.Name(),
			Language: language,
		}

		torrent.MagnetLink = e.ChildAttr("a[href^='magnet:']", "href")

		torrents = append(torrents, torrent)
	})

	formattedQuery := strings.ReplaceAll(query, " ", " ")
	searchURL, err := configs.BuildSearchURL("torrent9", formattedQuery, " ")

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

func GetMagnetLinkTorrent9(url string) (string, error) {
	cfg, _ := configs.GetScraperConfig("torrent9")
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
