package main

import (
	"command-dispatcher/internal/config"
	"command-dispatcher/internal/routes"
	"command-dispatcher/internal/subcribers"
	"sync"
)

func main() {
	config.Init()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		routes.Init()
	}()

	go func() {
		defer wg.Done()
		subcribers.Init()
	}()

	wg.Wait()

}
