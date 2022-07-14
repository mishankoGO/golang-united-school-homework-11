package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

// SafeCounter Mutex counter
type SafeCounter struct {
	mu    sync.Mutex
	users []user // users slice
}

// Inc method to increment users slice
func (c *SafeCounter) Inc(u user) {
	c.mu.Lock()
	c.users = append(c.users, u)
	c.mu.Unlock()
}

// getBatch func to load users in concurrent manner
func getBatch(n int64, pool int64) []user {

	// create wait group to monitor goroutines
	var wg sync.WaitGroup

	// Mutex counter
	res := SafeCounter{}

	// channel to limit number of concurrent load
	sem := make(chan struct{}, pool)

	// add goroutines to wait group in cycle
	for i := 0; int64(i) < n; i++ {
		wg.Add(1)

		// fill empty struct to channel
		sem <- struct{}{}

		// run goroutine
		go func(i int64) {
			u := getOne(i)
			<-sem
			res.Inc(u)
			wg.Done()
		}(int64(i)) // use index to avoid data Race
	}
	wg.Wait()
	return res.users
}
