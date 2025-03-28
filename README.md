# Torrent Scraper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A concurrent torrent search engine that aggregates results from multiple torrent sites.

## 🛠️ Initial Setup Required
**You must configure `configs/scrapers.json` before first use** - the template contains placeholder URLs that need replacement with current working domains.

## Features

- **Multi-site search**: Simultaneously queries 1337x, RARBG, NyaaSI, and Torrent9
- **Advanced filtering**: Filter by language, seeders count, and quality
- **Magnet links**: Automatic retrieval of magnet links for top results
- **JSON API**: Ready for integration with other applications
- **Configurable**: External configuration for sites and domains

## Installation

```bash
git clone https://github.com/yourusername/torrent-scraper.git
cd torrent-scraper
go build -o torrent-scraper cmd/main.go
```

## Configuration
Edit configs/scrapers.json (replace all placeholder values):
```json
{
  "sites": {
    "1337x": {
      "base_url": "https://www.url.domain",
      "search_path": "/search/{query}/1/",
      "domains": ["url.domain", "www.url.domain"]
    },
    "rarbg": {
      "base_url": "https://www.url.domain",
      "search_path": "/search/?search={query}",
      "domains": ["url.domain", "www.url.domain"]
    },
  }
}
```

## Usage
```bash
# Basic search (terminal output)
./torrent-scraper -q "Interstellar"

# JSON output with English results and magnets
./torrent-scraper -q "Dune" -lang english -m -limit 5 -o json
```

## Available Flags
```
Flag	Description	Example
-q	Search query	-q "Inception"
-o	Output format	-o json
-m	Fetch magnets	-m
-limit	Result limit	-limit 10
-lang	Language filter	-lang french
```

## Project Structure
```bash
configs/
└── scrapers.json       # Site configurations (REQUIRED SETUP)
internal/
├── core/              # Search logic
├── models/            # Data structures
└── scrapers/          # Site implementations
    ├── x1337.go
    ├── rarbg.go
    └── torrent9.go
```

License
MIT License - See LICENSE for details.