package main

import (
	"command-dispatcher/internal/config"
	"command-dispatcher/internal/routes"
	"command-dispatcher/internal/subcribers"
	"command-dispatcher/internal/worker"
	"log"
	"sync"
)

func main() {
	config.Init() // Initialize configuration
	// Start Asynq server in a goroutine
	var wg sync.WaitGroup
	wg.Add(3) // Now waiting for 3 goroutines: HTTP server, MQTT subscribers, and Asynq server

	go func() {
		defer wg.Done()
		log.Println("Starting HTTP routes...")
		routes.Init() // Initialize and start HTTP routes
	}()

	go func() {
		defer wg.Done()
		log.Println("Starting MQTT subscribers...")
		subcribers.Init() // Initialize and start MQTT subscribers
	}()

	go func() {
		defer wg.Done()
		log.Println("Starting queue worker...")
		worker.Init() // Initialize and start MQTT subscribers
	}()
	wg.Wait() // Wait for all goroutines to finish (though typically HTTP and Asynq servers run indefinitely)

	log.Println("Application stopped.")
}
