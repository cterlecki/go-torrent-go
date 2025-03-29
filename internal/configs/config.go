package configs

import (
	"fmt"
	"github.com/cterlecki/go-torrent-go/internal/models"
)

func PrintAvailableLanguages() {
	fmt.Println("Available languages:")
	for lang := range models.LanguagePatterns {
		fmt.Println("  -", lang)
	}
	fmt.Println("  - all (no filter)")
}

func ValidateLanguage(language string) error {
	if language != "" && language != "all" {
		valid := false
		for lang := range models.LanguagePatterns {
			if lang == language {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("Invalid language: %s\nUse -list-languages to see available options", language)
		}
	}
	return nil
}
