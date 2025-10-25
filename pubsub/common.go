package main

import (
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/dict"
)

type Request struct {
	name  string
	start time.Time
	end   time.Time
}

func NewRequest(name string) *Request {
	return &Request{
		name:  name,
		start: clock.TimeNow(),
	}
}

func (rq *Request) End() *Request {
	rq.end = clock.TimeNow()
	return rq
}

type ToggleParams struct {
	code     string
	isActive bool
}

func subscribe[T any](broker *PubSub[T], topic string, bufferSize int, lineMap dict.BoolMap) *Line[T] {
	line := broker.Subscribe(topic, bufferSize)
	lineMap[line.Name] = true
	return line
}

func runOrClose[T any](task func(T), data T, ok bool, line *Line[T], lineMap dict.BoolMap) {
	if ok {
		task(data)
	} else {
		lineMap[line.Name] = false
	}
}
