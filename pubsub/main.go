package main

import (
	"fmt"
	"sync"
	"time"
)

// Main application - run services and simulate requests

func main() {
	// For graceful shutdown, use WaitGroup
	// Otherwise, just run each service as a goroutine
	// go runCacheService()
	// go runMailService()
	// go runLogService()

	var wg sync.WaitGroup
	wg.Go(runCacheService)
	wg.Go(runMailService)
	wg.Go(runLogService)

	// Simulate requests
	go func() {
		addAccount(&Account{"roi", "Roi", "abc123"})
		addAccount(&Account{"john", "John", "def456"})
		time.Sleep(5 * time.Second)
		editAccount(&Account{"roi", "Roy", "def666"})
		time.Sleep(3 * time.Second)
		toggleAccount(&ToggleParams{"john", false})
		time.Sleep(5 * time.Second)
		toggleAccount(&ToggleParams{"john", true})
		time.Sleep(4 * time.Second)

		// Only for this example, so we can gracefully shutdown
		AccountBroker.Close()
		RequestBroker.Close()
		ToggleBroker.Close()
	}()

	// Normally, the web server here will prevent the
	// all goroutines are asleep - deadlock problem
	// server.Run()

	// For graceful shutdown, use WaitGroup
	wg.Wait()
	fmt.Println("Exiting program...")
}
