package main

import (
	"fmt"
	"time"

	"github.com/roidaradal/fn/clock"
)

func main() {
	TestFan()
}

func run(task func()) {
	start := clock.TimeNow()
	task()
	fmt.Printf("Time: %v\n\n", time.Since(start))
}
