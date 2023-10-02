package scraping_test

import (
	"context"
	"errors"
	"testing"

	"github.com/eriknylander/scraper/internal/filesystem"
	"github.com/eriknylander/scraper/internal/http"
	"github.com/eriknylander/scraper/internal/links"
	"github.com/eriknylander/scraper/internal/mock"
	"github.com/eriknylander/scraper/internal/scraping"
	"github.com/stretchr/testify/require"
)

func fetcherFailsAfterNMock(n int) func(ctx context.Context, url string) ([]byte, error) {
	return func(ctx context.Context, url string) ([]byte, error) {
		defer func() { n-- }()

		if n == 0 {
			return nil, errors.New("error")
		}

		return nil, nil
	}
}

func nestedLinksParserMock(l map[int][]string) func(b []byte) ([]string, error) {
	counter := 0
	return func(b []byte) ([]string, error) {
		defer func() { counter++ }()

		return l[counter], nil
	}
}

func TestScraper_Scrape(t *testing.T) {
	tests := []struct {
		name          string
		fetcher       http.Fetcher
		fileWriter    filesystem.FileWriter
		linkParser    links.Parser
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "fetch failed",
			fetcher:       &mock.FetcherMock{FetchFunc: func(ctx context.Context, url string) ([]byte, error) { return nil, errors.New("") }},
			fileWriter:    nil,
			linkParser:    nil,
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:          "file write failed",
			fetcher:       &mock.FetcherMock{FetchFunc: func(ctx context.Context, url string) ([]byte, error) { return []byte{}, nil }},
			fileWriter:    &mock.FileWriterMock{WriteFileFunc: func(path string, data []byte) error { return errors.New("error") }},
			linkParser:    nil,
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:       "link parser error",
			fetcher:    &mock.FetcherMock{FetchFunc: func(ctx context.Context, url string) ([]byte, error) { return []byte{}, nil }},
			fileWriter: &mock.FileWriterMock{WriteFileFunc: func(path string, data []byte) error { return nil }},
			linkParser: &mock.LinksParserMock{ParseLinksFunc: func(b []byte) ([]string, error) {
				return nil, errors.New("error")
			}},
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:       "not an HTML file",
			fetcher:    &mock.FetcherMock{FetchFunc: func(ctx context.Context, url string) ([]byte, error) { return []byte{}, nil }},
			fileWriter: &mock.FileWriterMock{WriteFileFunc: func(path string, data []byte) error { return nil }},
			linkParser: &mock.LinksParserMock{ParseLinksFunc: func(b []byte) ([]string, error) {
				return []string{"media/image.jpg"}, nil
			}},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "fetcher fails on nested page",
			fetcher: &mock.FetcherMock{
				FetchFunc: fetcherFailsAfterNMock(5),
			},
			fileWriter: &mock.FileWriterMock{WriteFileFunc: func(path string, data []byte) error { return nil }},
			linkParser: &mock.LinksParserMock{ParseLinksFunc: func(b []byte) ([]string, error) {
				return []string{"page1.html", "page2.html", "page3.html", "page5.html", "page6.html"}, nil
			}},
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:       "Successful nested",
			fetcher:    &mock.FetcherMock{FetchFunc: func(ctx context.Context, url string) ([]byte, error) { return nil, nil }},
			fileWriter: &mock.FileWriterMock{WriteFileFunc: func(path string, data []byte) error { return nil }},
			linkParser: &mock.LinksParserMock{
				ParseLinksFunc: nestedLinksParserMock(map[int][]string{
					0: {"page1.html", "page2.html", "page3.html"},
					1: {"sub1.html", "sub2.html", "style.css"},
					2: {"image.jpg"},
				}),
			},
			expectedCount: 8, // index + 7 assets
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scraping.NewScraper("", "", tt.fetcher, tt.linkParser, tt.fileWriter)

			res, err := s.Scrape(context.TODO())
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, res)
		})
	}
}
