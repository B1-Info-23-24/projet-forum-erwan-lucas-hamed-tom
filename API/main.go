package main

import (
	forum "forumApi/function"
	forumApi "forumApi/function"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	forumApi.InitDB()

	wg.Add(2)

	go func() {
		defer wg.Done()
		forum.StartAPIServer()

	}()

	wg.Wait()
}
