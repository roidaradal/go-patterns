package main

import (
	"fmt"
	"time"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/list"
)

type Data[T any] struct {
	index int
	item  T
}

type TransformFn[X any, Y any] = func(X) Y
type PipeFn[X any, Y any] = func(<-chan Data[X]) <-chan Data[Y]

var delay = 100 * time.Millisecond

func square(x int) int {
	time.Sleep(delay)
	y := x * x
	fmt.Printf("Square(%d) = %d\n", x, y)
	return y
}

func double(x int) int {
	time.Sleep(delay)
	y := 2 * x
	fmt.Printf("Double(%d) = %d\n", x, y)
	return y
}

func increment(x int) int {
	time.Sleep(delay)
	y := x + 1
	fmt.Printf("Inc(%d) = %d\n", x, y)
	return y
}

func Generate[T any](items ...T) <-chan Data[T] {
	outputCh := make(chan Data[T])
	go func() {
		for i, item := range items {
			outputCh <- Data[T]{i, item}
		}
		close(outputCh)
	}()
	return outputCh
}

func Consume[T any](channel <-chan Data[T], size int) []T {
	output := make([]T, size)
	for data := range channel {
		output[data.index] = data.item
	}
	return output
}

func Pipe[X any, Y any](fn TransformFn[X, Y]) PipeFn[X, Y] {
	return func(inputCh <-chan Data[X]) <-chan Data[Y] {
		outputCh := make(chan Data[Y])
		go func() {
			for input := range inputCh {
				outputCh <- Data[Y]{input.index, fn(input.item)}
			}
			close(outputCh)
		}()
		return outputCh
	}
}

func TestPipeline() {
	data := list.NumRange(1, 11)

	run(func() {
		out := data
		out = fn.Map(out, square)
		out = fn.Map(out, double)
		out = fn.Map(out, increment)
		fmt.Println(out)
	})

	run(func() {
		in := Generate(data...)
		out1 := Pipe(square)(in)
		out2 := Pipe(double)(out1)
		out3 := Pipe(increment)(out2)
		out := Consume(out3, len(data))
		fmt.Println(out)
	})
}
