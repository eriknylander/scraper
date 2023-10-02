package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/eriknylander/scraper/internal/filesystem"
	"github.com/eriknylander/scraper/internal/http"
	"github.com/eriknylander/scraper/internal/links"
	"github.com/eriknylander/scraper/internal/scraping"
)

var (
	siteToScrape    = flag.String("s", "", "The site to scrape")
	outputDirectory = flag.String("o", "", "Output directory")
)

func main() {
	ctx := context.Background()

	flag.Parse()

	if siteToScrape == nil || *siteToScrape == "" {
		fmt.Println("Missing site to scrape argument")
		flag.PrintDefaults()

		return
	}

	if outputDirectory == nil || *outputDirectory == "" {
		fmt.Println("Missing output directory argument, for usage, see scrape -h")
		flag.PrintDefaults()

		return
	}

	s := scraping.NewScraper(*siteToScrape, *outputDirectory, http.NewFetcher(), links.NewHTMLParser(), filesystem.NewFileWriter())

	downloaded, err := s.Scrape(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished scraping, downloaded: %d assets\n", downloaded)
}
