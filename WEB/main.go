package main

import (
	forumWeb "forumWeb/function"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		forumWeb.StartWebServer()
	}()

	wg.Wait()
}
