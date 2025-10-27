package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/roidaradal/fn/dict"
)

type DataFn[I any, O any] = func(I) (O, error)

type Input[T any] struct {
	index int
	item  T
}

type Output[T any] struct {
	index int
	item  T
	err   error
}

type Result[I any, O any] struct {
	success int
	output  map[int]O
	errors  map[int]error
}

func (r *Result[I, O]) Display(items []I) {
	fmt.Println("\nSuccess:", r.success)
	for i, item := range items {
		if dict.NoKey(r.output, i) {
			continue
		}
		fmt.Printf("In: %v Out: %v\n", item, r.output[i])
	}
	fmt.Println("Fail:", len(r.errors))
	for i, item := range items {
		if dict.NoKey(r.errors, i) {
			continue
		}
		fmt.Printf("In: %v Err: %s\n", item, r.errors[i].Error())
	}
	fmt.Println()
}

func NewResult[I any, O any]() *Result[I, O] {
	return &Result[I, O]{
		success: 0,
		output:  make(map[int]O),
		errors:  make(map[int]error),
	}
}

func Square(x int) (int, error) {
	if x == 3 || x == 6 {
		return 0, fmt.Errorf("cannot square %d", x)
	}
	time.Sleep(1 * time.Second) // artificial delay
	sq := x * x
	fmt.Printf("Square(%d) = %d\n", x, sq)
	return sq, nil
}

func LinearWorkers[I any, O any](items []I, fn DataFn[I, O]) *Result[I, O] {
	result := NewResult[I, O]()
	for i, item := range items {
		out, err := fn(item)
		if err == nil {
			result.success += 1
			result.output[i] = out
		} else {
			result.errors[i] = err
		}
	}
	return result
}

func ConcurrentWorkers[I any, O any](items []I, fn DataFn[I, O], numWorkers int) *Result[I, O] {
	// Input and output channels
	inputCh := make(chan Input[I])
	outputCh := make(chan Output[O], numWorkers) // buffered, otherwise deadlocks

	// Worker function
	worker := func(id int, inputCh <-chan Input[I], outputCh chan<- Output[O]) {
		count := 0
		for input := range inputCh {
			out, err := fn(input.item)
			outputCh <- Output[O]{input.index, out, err}
			count += 1
		}
		fmt.Printf("Worker %d did %d jobs\n", id, count)
	}

	// Spawn the workers
	var wg sync.WaitGroup
	for id := range numWorkers {
		wg.Go(func() {
			worker(id, inputCh, outputCh)
		})
	}

	// Feed the input data to input channel
	go func() {
		for i, item := range items {
			inputCh <- Input[I]{i, item}
		}
		close(inputCh)
	}()

	// Wait for all workers to finish, close the output channel
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	// Get the results
	result := NewResult[I, O]()
	for out := range outputCh {
		if out.err == nil {
			result.success += 1
			result.output[out.index] = out.item
		} else {
			result.errors[out.index] = out.err
		}
	}
	return result
}

func TestPool() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8}

	run(func() {
		fmt.Println("Linear Workers")
		result := LinearWorkers(data, Square)
		result.Display(data)
	})

	run(func() {
		fmt.Println("Concurrent Workers")
		result := ConcurrentWorkers(data, Square, 4)
		result.Display(data)
	})
}
