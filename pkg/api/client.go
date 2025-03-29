package api

import (
	"fmt"
	"github.com/cterlecki/go-torrent-go/internal/configs"
	"github.com/cterlecki/go-torrent-go/internal/core"
	"github.com/cterlecki/go-torrent-go/internal/models"
)

type TorrentClient struct {
	initialized bool
	configPath  string
}

func New(configPath string) (*TorrentClient, error) {
	client := &TorrentClient{configPath: configPath}
	if err := client.Initialize(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *TorrentClient) Initialize() error {
	if err := configs.LoadScrapersConfig(c.configPath); err != nil {
		return err
	}
	c.initialized = true
	return nil
}

func (c *TorrentClient) Search(query, language string, fetchMagnets bool, limit int) ([]models.Torrent, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	results := core.SearchAllSites(query, fetchMagnets, language)
	if limit > 0 && limit < len(results) {
		return results[:limit], nil
	}
	return results, nil
}
