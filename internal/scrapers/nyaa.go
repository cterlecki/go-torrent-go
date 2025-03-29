package scrapers

import (
	"fmt"
	"github.com/cterlecki/go-torrent-go/internal/configs"
	"github.com/cterlecki/go-torrent-go/internal/models"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
	"time"
)

type NyaaScraper struct{}

func (s NyaaScraper) Name() string {
	return "NyaaSI"
}

func (s NyaaScraper) Search(query string) ([]models.Torrent, error) {
	cfg, _ := configs.GetScraperConfig("nyaa")
	var torrents []models.Torrent

	c := colly.NewCollector(
		colly.AllowedDomains(cfg.Domains...),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	c.OnHTML("table.torrent-list tbody tr", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("td:nth-child(2) a:last-child"))

		quality := models.ExtractQuality(name)
		language := models.DetectLanguage(name)

		if strings.Contains(strings.ToLower(name), "raw") {
			language = "japanese"
		}

		torrent := models.Torrent{
			Name:       name,
			URL:        cfg.BaseURL + e.ChildAttr("td:nth-child(2) a:last-child", "href"),
			Size:       e.ChildText("td:nth-child(4)"),
			UploadDate: e.ChildText("td:nth-child(5)"),
			Quality:    quality,
			Source:     s.Name(),
			Language:   language,
			MagnetLink: e.ChildAttr("td:nth-child(3) a:nth-child(2)", "href"),
		}

		seedersStr := e.ChildText("td:nth-child(6)")
		seeders, err := strconv.Atoi(seedersStr)
		if err == nil {
			torrent.Seeders = seeders
		}

		leechersStr := e.ChildText("td:nth-child(7)")
		leechers, err := strconv.Atoi(leechersStr)
		if err == nil {
			torrent.Leechers = leechers
		}

		torrents = append(torrents, torrent)
	})

	formattedQuery := strings.ReplaceAll(query, " ", "+")
	searchURL, err := configs.BuildSearchURL("nyaa", formattedQuery, "+")

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
