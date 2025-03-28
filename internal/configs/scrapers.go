package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ScraperConfig struct {
	BaseURL    string   `json:"base_url"`
	SearchPath string   `json:"search_path"`
	Domains    []string `json:"domains"`
}

type ScrapersConfig struct {
	Sites map[string]ScraperConfig `json:"sites"`
}

var scraperConfig *ScrapersConfig

func LoadScrapersConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	scraperConfig = &ScrapersConfig{}
	return json.Unmarshal(file, scraperConfig)
}

func GetScraperConfig(name string) (ScraperConfig, bool) {
	cfg, exists := scraperConfig.Sites[name]
	return cfg, exists
}

func BuildSearchURL(siteName, query string, replacement string) (string, error) {
	cfg, exists := GetScraperConfig(siteName)
	if !exists {
		return "", fmt.Errorf("scraper %s not configured", siteName)
	}

	formattedQuery := strings.ReplaceAll(query, " ", replacement)
	searchURL := cfg.BaseURL + strings.Replace(cfg.SearchPath, "{query}", formattedQuery, 1)
	return searchURL, nil
}
