package scraping

import "sync"

// claimer keeps track of claimed pages.
type claimer struct {
	*sync.Mutex
	visited map[string]bool
}

func newClaimer() claimer {
	return claimer{
		Mutex:   &sync.Mutex{},
		visited: make(map[string]bool),
	}
}

// claimURL takes a url and checks whether is has already been claimed or not. Returns false if the url has already been claimed, claims and returns true if the url has not been claimed.
func (c claimer) claimURL(url string) bool {
	c.Lock()
	defer c.Unlock()

	if ok := c.visited[url]; ok {
		return false
	}

	c.visited[url] = true

	return true
}
