package main

import (
	"fmt"
	"time"

	"github.com/roidaradal/fn/clock"
)

var runStart time.Time

func main() {
	// TestTasks()
	// TestActions()
	// TestSimpleData()
	// TestData()
	TestRequests()
}

func run(task func()) {
	start := clock.TimeNow()
	task()
	fmt.Printf("Time: %v\n\n", time.Since(start))
}

func elapsed() string {
	duration := time.Since(runStart).Seconds()
	return fmt.Sprintf("%.0f", duration)
}
