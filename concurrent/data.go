package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type SimpleDataFn[I any, O any] = func(I) O
type DataFn[I any, O any] = func(I) (O, error)

type Data[T any] struct {
	index  int
	output T
}

func Square(x int) int {
	time.Sleep(1 * time.Second) // artificial delay
	sq := x * x
	fmt.Printf("Square(%d) = %d\n", x, sq)
	return sq
}

func SquareErr(x int) (int, error) {
	if x == 4 {
		return 0, errors.New("cannot square 4")
	}
	return Square(x), nil
}

func LinearSimpleData[I any, O any](items []I, fn SimpleDataFn[I, O]) []O {
	output := make([]O, len(items))
	for i, item := range items {
		output[i] = fn(item)
	}
	return output
}

func LinearData[I any, O any](items []I, fn DataFn[I, O]) ([]O, error) {
	output := make([]O, len(items))
	for i, item := range items {
		out, err := fn(item)
		if err != nil {
			return output, err
		}
		output[i] = out
	}
	return output, nil
}

func ConcurrentSimpleData[I any, O any](items []I, fn SimpleDataFn[I, O]) []O {
	var wg sync.WaitGroup
	result := make(chan Data[O])
	for i, item := range items {
		wg.Go(func() {
			result <- Data[O]{i, fn(item)}
		})
	}

	// Wait for everything to finish, then close result channel
	// Wrap in goroutine so that it's non-blocking
	go func() {
		wg.Wait()
		close(result)
	}()

	// Receive data from result channel
	output := make([]O, len(items))
	for data := range result {
		output[data.index] = data.output
	}
	return output
}

func ConcurrentData[I any, O any](items []I, fn DataFn[I, O]) ([]O, error) {
	var eg errgroup.Group
	result := make(chan Data[O])
	for i, item := range items {
		eg.Go(func() error {
			out, err := fn(item)
			result <- Data[O]{i, out}
			return err
		})
	}

	var finalErr error
	go func() {
		finalErr = eg.Wait()
		close(result)
	}()

	output := make([]O, len(items))
	for data := range result {
		output[data.index] = data.output
	}
	return output, finalErr
}

func TestSimpleData() {
	data := []int{1, 2, 3, 4, 5}

	run(func() {
		fmt.Println("Linear Data")
		output := LinearSimpleData(data, Square)
		fmt.Println("Out:", output)
	})

	run(func() {
		fmt.Println("Concurrent Data")
		output := ConcurrentSimpleData(data, Square)
		fmt.Println("Out:", output)
	})
}

func TestData() {
	data := []int{1, 2, 3, 4, 5}

	run(func() {
		fmt.Println("Linear Data")
		output, err := LinearData(data, SquareErr)
		fmt.Println("Out:", output)
		fmt.Println("Err:", err)
	})

	run(func() {
		fmt.Println("Concurrent Data")
		output, err := ConcurrentData(data, SquareErr)
		fmt.Println("Out:", output)
		fmt.Println("Err:", err)
	})
}
