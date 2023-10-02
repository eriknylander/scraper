package http

import (
	"context"
	"io"
	"net/http"
)

// Fetcher describes an interface for fetching websites.
type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// fetcher implements the Fetcher-interface
type fetcher struct{}

func NewFetcher() Fetcher {
	return fetcher{}
}

// Fetcher fetches a website and returns the body as a byte slice.
func (f fetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
