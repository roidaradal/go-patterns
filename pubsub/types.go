package main

import (
	"time"

	"github.com/roidaradal/fn/clock"
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
