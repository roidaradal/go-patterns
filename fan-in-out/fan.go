package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/roidaradal/fn/conv"
	"github.com/roidaradal/fn/list"
)

type Task[X any, Y any] = func(X) Y

type Output[T any] struct {
	index int
	item  T
}

func expand(n int) int {
	time.Sleep(1 * time.Second)
	expanded := fmt.Sprintf("%d%d", n, n)
	fmt.Printf("Expand %d => %s\n", n, expanded)
	return conv.ParseInt(expanded)
}

func Linear[X any, Y any](items []X, task Task[X, Y]) []Y {
	results := make([]Y, len(items))
	for i, item := range items {
		results[i] = task(item)
	}
	return results
}

func FanOutIn[X any, Y any](items []X, task Task[X, Y], numWorkers int) []Y {
	channels := FanOut(items, task, numWorkers)
	resultCh := FanIn(numWorkers, channels...)

	results := make([]Y, len(items))
	for out := range resultCh {
		results[out.index] = out.item
	}
	return results
}

func FanOut[X any, Y any](items []X, task Task[X, Y], numWorkers int) []<-chan Output[Y] {
	channels := make([]<-chan Output[Y], numWorkers)
	numItems := len(items)

	// Fan-out the workload
	for workerID := range numWorkers {
		workerCh := make(chan Output[Y])
		channels[workerID] = workerCh

		// start worker goroutine
		go func() {
			count := 0
			for j := workerID; j < numItems; j += numWorkers {
				out := task(items[j])
				workerCh <- Output[Y]{j, out}
				count += 1
			}
			fmt.Printf("Worker %d finished %d tasks\n", workerID, count)
			close(workerCh)
		}()
	}
	return channels
}

func FanIn[T any](numWorkers int, channels ...<-chan Output[T]) <-chan Output[T] {
	result := make(chan Output[T], numWorkers)

	var wg sync.WaitGroup
	for _, workerCh := range channels {
		wg.Go(func() {
			for out := range workerCh {
				result <- out
			}
		})
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func TestFan() {
	data := list.NumRange(1, 21)

	run(func() {
		fmt.Println("Linear")
		results := Linear(data, expand)
		fmt.Println(len(results), results)
	})

	run(func() {
		fmt.Println("Fan-Out/Fan-In")
		results := FanOutIn(data, expand, 4)
		fmt.Println(len(results), results)
	})
}
