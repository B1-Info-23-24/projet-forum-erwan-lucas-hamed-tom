package main

import (
	forum "forum/function"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		forum.StartWebServer()
	}()

	go func() {
		defer wg.Done()
		forum.StartAPIServer()
	}()

	wg.Wait()
}
