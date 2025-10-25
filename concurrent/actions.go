package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type ActionFn = func() error
type CtxActionFn = func(context.Context) error

func LinearActions(actions []ActionFn) error {
	runStart = time.Now()
	for _, action := range actions {
		if err := action(); err != nil {
			return err
		}
	}
	return nil
}

// Runs all actions concurrently, returns first error (if any)
// Once action has been started, it finishes its task to the end (no cancellation)
func ConcurrentActions(actions []ActionFn) error {
	var eg errgroup.Group
	runStart = time.Now()
	for _, action := range actions {
		eg.Go(action)
	}
	return eg.Wait()
}

// Runs context actions concurrently, return first error (if any)
// Context is passed to the functions (and usually to children functions),
// so that we can check if the context has been cancelled and end early
func ConcurrentCtxActions(ctxActions []CtxActionFn, timeout float64) error {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		duration := time.Duration(timeout) * time.Second
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}

	group, ctx := errgroup.WithContext(ctx)
	runStart = time.Now()
	for _, ctxAction := range ctxActions {
		group.Go(func() error {
			return ctxAction(ctx)
		})
	}
	return group.Wait()
}

func newAction(duration int) ActionFn {
	return func() error {
		fmt.Printf("[%s] Task %d started\n", elapsed(), duration)
		time.Sleep(time.Duration(duration) * time.Second)

		if duration == 6 || duration == 10 {
			fmt.Printf("[%s] Task %d fail\n", elapsed(), duration)
			return fmt.Errorf("bad input: %d", duration)
		}

		fmt.Printf("[%s] Task %d done\n", elapsed(), duration)
		return nil
	}
}

func newCtxAction(duration int) CtxActionFn {
	return func(ctx context.Context) error {
		idle := duration / 2
		work := duration - idle

		// Artificial pause before starting
		time.Sleep(time.Duration(idle) * time.Second)

		// Check if context has been cancelled
		// Usually put this in front of an expensive task:
		// Check first if context cancelled to avoid spending resources doing expensive task
		select {
		default:
		case <-ctx.Done():
			fmt.Printf("[%s] Task %d cancelled\n", elapsed(), duration)
			return ctx.Err()
		}

		fmt.Printf("[%s] Task %d started\n", elapsed(), duration)

		if duration == 6 || duration == 10 {
			fmt.Printf("[%s] Task %d fail\n", elapsed(), duration)
			return fmt.Errorf("bad input: %d", duration)
		}

		time.Sleep(time.Duration(work) * time.Second)
		fmt.Printf("[%s] Task %d done\n", elapsed(), duration)
		return nil
	}

}

func TestActions() {
	actions := []ActionFn{
		newAction(2),
		newAction(4),
		newAction(6),
		newAction(8),
		newAction(10),
	}

	run(func() {
		fmt.Println("Linear Actions")
		err := LinearActions(actions)
		fmt.Println("Error:", err)
	})

	ctxActions := []CtxActionFn{
		newCtxAction(2),
		newCtxAction(4),
		newCtxAction(6),
		newCtxAction(8),
		newCtxAction(10),
	}

	run(func() {
		fmt.Println("Concurrent Actions")
		// err := ConcurrentActions(actions)
		// err := ConcurrentCtxActions(ctxActions, 0)
		err := ConcurrentCtxActions(ctxActions, 2.5)
		fmt.Println("Error:", err)
	})
}
