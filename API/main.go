package main

import (
	forum "forumAPI/function"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		forum.StartAPIServer()
	}()
	wg.Wait()
}
