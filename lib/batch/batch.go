package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

type SafeCounter struct {
	mu sync.Mutex
	v  []user
}

func (c *SafeCounter) Inc(u user) {
	c.mu.Lock()
	c.v = append(c.v, u)
	c.mu.Unlock()
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) []user {

	var wg sync.WaitGroup
	res := SafeCounter{}

	sem := make(chan struct{}, pool)

	for i := 0; int64(i) < n; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int64) {
			u := getOne(i)
			<-sem
			res.Inc(u)
			wg.Done()
		}(int64(i))
	}
	wg.Wait()
	return res.v
}
