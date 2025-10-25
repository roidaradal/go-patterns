package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type RequestFn = func(*Request) error

type Request struct {
	logs []string
	mu   sync.RWMutex
}

func NewRequest() *Request {
	return &Request{logs: make([]string, 0)}
}

func (rq *Request) SubRequest() *Request {
	return &Request{logs: make([]string, 0)}
}

func (rq *Request) MergeLogs(srq *Request) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.logs = append(rq.logs, srq.logs...)
}

func (rq *Request) AddLog(message string) {
	rq.logs = append(rq.logs, message)
}

func (rq *Request) Output() string {
	return strings.Join(rq.logs, "\n")
}

func LinearRequests(rq *Request, requests []RequestFn) error {
	runStart = time.Now()
	for _, request := range requests {
		if err := request(rq); err != nil {
			return err
		}
	}
	return nil
}

func ConcurrentRequests(rq *Request, requests []RequestFn) error {
	var eg errgroup.Group
	runStart = time.Now()
	for _, request := range requests {
		eg.Go(func() error {
			srq := rq.SubRequest()
			err := request(srq)
			rq.MergeLogs(srq)
			return err
		})
	}
	return eg.Wait()
}

func newRequestFn(duration int) RequestFn {
	return func(rq *Request) error {
		rq.AddLog(fmt.Sprintf("[%s] Req %d started", elapsed(), duration))
		time.Sleep(time.Duration(duration) * time.Second)

		if duration == 3 || duration == 5 {
			rq.AddLog(fmt.Sprintf("[%s] Req %d fail", elapsed(), duration))
			return fmt.Errorf("bad input: %d", duration)
		}

		rq.AddLog(fmt.Sprintf("[%s] Req %d done", elapsed(), duration))
		return nil
	}
}

func TestRequests() {
	requests := []RequestFn{
		newRequestFn(1),
		newRequestFn(2),
		newRequestFn(3),
		newRequestFn(4),
		newRequestFn(5),
	}

	run(func() {
		fmt.Println("Linear Requests")
		rq := NewRequest()
		err := LinearRequests(rq, requests)
		fmt.Println(rq.Output())
		fmt.Println("Error:", err)
	})

	run(func() {
		fmt.Println("Concurrent Requests")
		rq := NewRequest()
		err := ConcurrentRequests(rq, requests)
		fmt.Println(rq.Output())
		fmt.Println("Error:", err)
	})
}
