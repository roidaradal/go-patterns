package main

import (
	"fmt"
	"sync"
	"time"
)

type TaskFn = func()

func LinearTasks(tasks []TaskFn) {
	for _, task := range tasks {
		task()
	}
}

func ConcurrentTasks(tasks []TaskFn) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Go(task)
	}
	wg.Wait()
}

func newTask(duration int) TaskFn {
	return func() {
		time.Sleep(time.Duration(duration) * time.Second)
		fmt.Printf("Task %d done\n", duration)
	}
}

func TestTasks() {
	tasks := []TaskFn{
		newTask(1),
		newTask(2),
		newTask(3),
		newTask(4),
	}

	run(func() {
		fmt.Println("Linear Tasks")
		LinearTasks(tasks)
	})

	run(func() {
		fmt.Println("Concurrent Tasks")
		ConcurrentTasks(tasks)
	})
}
