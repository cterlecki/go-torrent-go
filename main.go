package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cterlecki/go-torrent-go/internal/configs"
	"github.com/cterlecki/go-torrent-go/internal/core"
	"github.com/cterlecki/go-torrent-go/internal/models"
	"log"
	"os"
	"strings"
)

func main() {
	queryFlag := flag.String("q", "", "Search query")
	outputFlag := flag.String("o", "terminal", "Output format (terminal, json)")
	magnetsFlag := flag.Bool("m", false, "Fetch magnet links for top results")
	limitFlag := flag.Int("limit", 0, "Limit number of results (0 for no limit)")
	langFlag := flag.String("lang", "", "Filter by language")
	listLangFlag := flag.Bool("list-languages", false, "List all available languages")
	flag.Parse()

	err := configs.LoadScrapersConfig("./internal/configs/scrapers.json")
	if err != nil {
		log.Fatalf("Failed to load scrapers configs: %v", err)
	}

	if *listLangFlag {
		configs.PrintAvailableLanguages()
		os.Exit(0)
	}

	query := *queryFlag
	if query == "" {
		query = strings.Join(flag.Args(), " ")
	}

	if query == "" {
		fmt.Println("Usage: torrent-scraper -q \"search query\" [-o terminal|json] [-m] [-limit N] [-lang language]")
		fmt.Println("Use -list-languages to see available language options")
		os.Exit(1)
	}

	language := strings.ToLower(*langFlag)
	if err := configs.ValidateLanguage(language); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if language != "" && language != "all" {
		fmt.Printf(" (Language: %s)", language)
	}

	torrents := core.SearchAllSites(query, *magnetsFlag, language)

	if *limitFlag > 0 && *limitFlag < len(torrents) {
		torrents = torrents[:*limitFlag]
	}

	switch *outputFlag {
	case "json":
		outputJSON(torrents)
	default:
		outputTerminal(torrents)
	}
}

func outputJSON(torrents []models.Torrent) {
	jsonData, err := json.MarshalIndent(torrents, "", "  ")
	if err != nil {
		log.Fatal("Error generating JSON:", err)
	}
	fmt.Println(string(jsonData))
}

func outputTerminal(torrents []models.Torrent) {
	fmt.Printf("Found %d torrents:\n\n", len(torrents))
	for i, t := range torrents {
		fmt.Printf("%d. %s\n", i+1, t.Name)
		fmt.Printf("   Quality: %s | Size: %s | Language: %s\n",
			t.Quality, t.Size, t.Language)
		fmt.Printf("   Seeders: %d | Leechers: %d | Source: %s\n",
			t.Seeders, t.Leechers, t.Source)
		if t.UploadDate != "" {
			fmt.Printf("   Upload Date: %s\n", t.UploadDate)
		}
		if t.MagnetLink != "" {
			fmt.Printf("   Magnet: %s\n", t.MagnetLink)
		}
		fmt.Println()
	}
}
