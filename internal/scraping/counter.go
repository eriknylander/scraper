package scraping

import "sync"

type counter struct {
	*sync.Mutex
	counter int
}

func newCounter() counter {
	return counter{Mutex: &sync.Mutex{}, counter: 0}
}

func (c *counter) increase() {
	c.Lock()
	defer c.Unlock()

	c.counter += 1
}
