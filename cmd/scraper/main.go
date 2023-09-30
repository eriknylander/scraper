package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eriknylander/scraper/internal/links"
)

func main() {
	ctx := context.Background()
	b, err := getPage(ctx, "https://books.toscrape.com")
	if err != nil {
		panic(err)
	}

	l, err := links.NewHTMLParser().ParseLinks(b)
	if err != nil {
		panic(err)
	}

	for i := range l {
		fmt.Println(l[i])
	}
}

func getPage(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
