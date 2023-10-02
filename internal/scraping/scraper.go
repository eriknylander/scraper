package scraping

import (
	"context"
	"fmt"
	"path"
	"sync"

	"github.com/eriknylander/scraper/internal/filesystem"
	"github.com/eriknylander/scraper/internal/http"
	"github.com/eriknylander/scraper/internal/links"
)

type Scraper struct {
	claimer         claimer
	counter         counter
	baseURL         string
	outputDirectory string
	fetcher         http.Fetcher
	linksParser     links.Parser
	fileWriter      filesystem.FileWriter
}

func NewScraper(baseURL string, outputDirectory string, fetcher http.Fetcher, linksParser links.Parser, fileWriter filesystem.FileWriter) Scraper {
	return Scraper{claimer: newClaimer(), counter: newCounter(), baseURL: baseURL, outputDirectory: outputDirectory, fetcher: fetcher, linksParser: linksParser, fileWriter: fileWriter}
}

func (s *Scraper) Scrape(ctx context.Context) (int, error) {
	s.claimer.claimURL("")
	s.claimer.claimURL("index.html")
	s.claimer.claimURL("index.htm")

	if err := s.scrape(ctx, ""); err != nil {
		return 0, err
	}

	return s.counter.counter, nil
}

func (s *Scraper) scrape(ctx context.Context, url string) error {
	durl := fmt.Sprintf("%s/%s", s.baseURL, url)
	d, err := s.fetcher.Fetch(ctx, durl)
	if err != nil {
		return err
	}

	if url == "" {
		url = "index.html"
	}

	if err := s.fileWriter.WriteFile(path.Join(s.outputDirectory, url), d); err != nil {
		return err
	}

	fmt.Printf("Downloaded %s\n", durl)
	s.counter.increase()

	if !isHTMLPage(url) {
		return nil
	}

	links, err := s.linksParser.ParseLinks(d)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, l := range links {
		if isExternalLink(l) {
			continue
		}

		nl := normalizePath(url, l)
		if !s.claimer.claimURL(nl) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := s.scrape(ctx, nl); err != nil {
				errCh <- err
			}
		}()
	}

	finishedCh := make(chan struct{})
	go func() {
		wg.Wait()
		finishedCh <- struct{}{}
	}()

	select {
	case <-finishedCh:
		return nil
	case err := <-errCh:
		return err
	}
}
