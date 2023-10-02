package mock

import "context"

type FetcherMock struct {
	FetchFunc func(ctx context.Context, url string) ([]byte, error)
}

func (fm *FetcherMock) Fetch(ctx context.Context, url string) ([]byte, error) {
	return fm.FetchFunc(ctx, url)
}
