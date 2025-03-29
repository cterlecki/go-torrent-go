package scrapers

import (
	"github.com/cterlecki/go-torrent-go/internal/models"
	"strconv"
	"strings"
)

type TorrentScraper interface {
	Search(query string) ([]models.Torrent, error)
	Name() string
}

func parseInt(s string) int {
	num, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return num
}
