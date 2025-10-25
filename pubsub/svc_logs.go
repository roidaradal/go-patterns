package main

import "fmt"

// Log service : Subscriber example

func runLogService() {
	logRequest := RequestBroker.Subscribe(E_END_REQUEST, 10)

	fmt.Println("[Log] Service started...")
	for rq := range logRequest.Channel {
		duration := rq.end.Sub(rq.start)
		fmt.Printf("[Log] Request: %s, Duration: %v\n", rq.name, duration)
	}
	fmt.Println("[Log] Service stopped...")
}
