package models

import (
	"regexp"
	"strings"
)

type Torrent struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Seeders    int    `json:"seeders"`
	Leechers   int    `json:"leechers"`
	Size       string `json:"size"`
	UploadDate string `json:"upload_date,omitempty"`
	Quality    string `json:"quality"`
	Source     string `json:"source"`
	MagnetLink string `json:"magnet_link,omitempty"`
	Language   string `json:"language"`
}

var LanguagePatterns = map[string]string{
	"english":    `\beng\b|\benglish\b|\b(us|uk)[\.-]\b`,
	"french":     `\bfr(e|a)?\b|\bfrench\b|\bvf\b|\bvff\b`,
	"spanish":    `\besp\b|\bspa\b|\bspanish\b`,
	"german":     `\bger\b|\bdeu\b|\bgerman\b`,
	"italian":    `\bita?\b|\bitalian\b`,
	"portuguese": `\bpor\b|\bpt\b|\bportuguese\b|\bbrazil\b`,
	"russian":    `\brus\b|\brussian\b`,
	"japanese":   `\bjapanese\b|\bjpn\b|\bjap\b`,
	"korean":     `\bkorean\b|\bkor\b`,
	"chinese":    `\bchinese\b|\bchi\b|\bman?\b|\bcantonese\b`,
	"hindi":      `\bhindi\b`,
	"multi":      `\bmulti(language)?\b|\bmulti-sub\b`,
}

func DetectLanguage(name string) string {
	name = strings.ToLower(name)

	for lang, pattern := range LanguagePatterns {
		if match, _ := regexp.MatchString(pattern, name); match {
			return lang
		}
	}

	return "unknown"
}

func ExtractQuality(name string) string {
	name = strings.ToLower(name)

	resolutionPatterns := map[string]string{
		"4k":    `\b4k\b|\b2160p\b`,
		"1080p": `\b1080p\b`,
		"720p":  `\b720p\b`,
		"480p":  `\b480p\b`,
	}

	sourcePatterns := map[string]string{
		"bluray": `\bbluray\b|\bblu-ray\b`,
		"remux":  `\bremux\b`,
		"web-dl": `\bweb-?dl\b`,
		"webrip": `\bwebrip\b|\bweb-?rip\b`,
		"hdtv":   `\bhdtv\b`,
	}

	var quality string
	for res, pattern := range resolutionPatterns {
		match, _ := regexp.MatchString(pattern, name)
		if match {
			quality = res
			break
		}
	}

	for src, pattern := range sourcePatterns {
		match, _ := regexp.MatchString(pattern, name)
		if match {
			if quality != "" {
				quality += " "
			}
			quality += src
			break
		}
	}

	if quality == "" {
		quality = "unknown"
	}

	return quality
}
