package main

import (
	"context"
	"fmt"

	"github.com/eriknylander/scraper/internal/filesystem"
	"github.com/eriknylander/scraper/internal/http"
	"github.com/eriknylander/scraper/internal/links"
	"github.com/eriknylander/scraper/internal/scraping"
)

func main() {
	ctx := context.Background()

	s := scraping.NewScraper("http://books.toscrape.com", "C:\\Users\\nylan\\go\\src\\github.com\\eriknylander\\scraper\\downloaded", http.NewFetcher(), links.NewHTMLParser(), filesystem.NewFileWriter())

	downloaded, err := s.Scrape(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished scraping, downloaded: %d assets\n", downloaded)
}
