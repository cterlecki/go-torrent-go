package core

import (
	"log"
	"main/internal/models"
	"main/internal/scrapers"
	"sort"
	"sync"
)

func SearchAllSites(query string, fetchMagnets bool, language string) []models.Torrent {
	var wg sync.WaitGroup
	scrapersList := []scrapers.TorrentScraper{
		&scrapers.X1337Scraper{},
		&scrapers.RARBGScraper{},
		&scrapers.NyaaScraper{},
		&scrapers.Torrent9Scraper{},
	}

	resultChan := make(chan []models.Torrent, len(scrapersList))

	for _, scraper := range scrapersList {
		wg.Add(1)
		go func(s scrapers.TorrentScraper) {
			defer wg.Done()

			torrents, err := s.Search(query)
			if err != nil {
				log.Printf("Error searching %s: %v", s.Name(), err)
				resultChan <- []models.Torrent{}
				return
			}

			resultChan <- torrents
		}(scraper)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allTorrents []models.Torrent
	for torrents := range resultChan {
		allTorrents = append(allTorrents, torrents...)
	}

	if language != "" && language != "all" {
		var filteredTorrents []models.Torrent
		for _, t := range allTorrents {
			if t.Language == language ||
				(language == "english" && t.Language == "unknown") ||
				(t.Language == "multi") {
				filteredTorrents = append(filteredTorrents, t)
			}
		}
		allTorrents = filteredTorrents
	}

	sort.Slice(allTorrents, func(i, j int) bool {
		return allTorrents[i].Seeders > allTorrents[j].Seeders
	})

	if fetchMagnets && len(allTorrents) > 0 {
		limit := 20
		if len(allTorrents) < limit {
			limit = len(allTorrents)
		}

		var magnetWg sync.WaitGroup
		for i := 0; i < limit; i++ {
			if allTorrents[i].MagnetLink == "" {
				magnetWg.Add(1)
				go func(index int) {
					defer magnetWg.Done()
					switch allTorrents[index].Source {
					case "1337x":
						magnetLink, err := scrapers.GetMagnetLinkX1337(allTorrents[index].URL)
						if err == nil {
							allTorrents[index].MagnetLink = magnetLink
						}
					case "Torrent9":
						magnetLink, err := scrapers.GetMagnetLinkTorrent9(allTorrents[index].URL)
						if err == nil {
							allTorrents[index].MagnetLink = magnetLink
						}
					}
				}(i)
			}
		}
		magnetWg.Wait()
	}

	return allTorrents
}
