package main

import (
	forum "forumWeb/function"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		forum.StartWebServer()
	}()
	wg.Wait()
}
