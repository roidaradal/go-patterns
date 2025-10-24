package main

import "fmt"

func runLogService() {
	logRequestCh := RequestBroker.Subscribe(E_END_REQUEST, 10)

	fmt.Println("[LOG] Service started...")
	for rq := range logRequestCh {
		duration := rq.end.Sub(rq.start)
		fmt.Printf("[LOG] Request: %s, Duration: %v\n", rq.name, duration)
	}
}
