package main

import (
	"fmt"
	"sync"
)

// Message Broker module

// Event names
const (
	E_ADD_ACCOUNT    string = "add-account"
	E_TOGGLE_ACCOUNT string = "toggle-account"
	E_EDIT_ACCOUNT   string = "edit-account"
	E_END_REQUEST    string = "end-request"
)

// Brokers
var (
	AccountBroker = NewPubSub[*Account]("Account")
	RequestBroker = NewPubSub[*Request]("Request")
	ToggleBroker  = NewPubSub[*ToggleParams]("Toggle")
)

// Message Broker (can handle one type)

type Line[T any] struct {
	Name    string
	Channel <-chan T
}

type PubSub[T any] struct {
	mu          sync.RWMutex
	name        string
	subscribers map[string][]chan T // map: topic => list of channels to subscribers
}

func NewPubSub[T any](name string) *PubSub[T] {
	return &PubSub[T]{
		name:        name,
		subscribers: make(map[string][]chan T),
	}
}

func (ps *PubSub[T]) Subscribe(topic string, bufferSize int) *Line[T] {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	channel := make(chan T, bufferSize)
	ps.subscribers[topic] = append(ps.subscribers[topic], channel)

	subscriber := &Line[T]{
		Name:    fmt.Sprintf("%s.%s", ps.name, topic),
		Channel: channel,
	}
	return subscriber
}

func (ps *PubSub[T]) Publish(topic string, message T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, subscriber := range ps.subscribers[topic] {
		subscriber <- message
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, subscribers := range ps.subscribers {
		for _, channel := range subscribers {
			close(channel)
		}
		fmt.Println("[TOPIC] Closed:", topic)
	}
}
