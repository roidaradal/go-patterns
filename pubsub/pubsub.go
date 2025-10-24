package main

import "sync"

// Homogeneous Messages (one type)

type PubSub[T any] struct {
	mu          sync.RWMutex
	subscribers map[string][]chan T // map: topic => list of channels to subscribers
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[string][]chan T),
	}
}

func (ps *PubSub[T]) Subscribe(topic string, bufferSize int) <-chan T {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subscriber := make(chan T, bufferSize)
	ps.subscribers[topic] = append(ps.subscribers[topic], subscriber)
	return subscriber
}

func (ps *PubSub[T]) Publish(topic string, message T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, subscriber := range ps.subscribers[topic] {
		subscriber <- message
	}
}
