package main

import (
	"fmt"
	"time"
)

type Server struct {
	host    string
	port    int
	timeout time.Duration
	useTLS  bool
}

type OptFunc[T any] func(*T)

func WithHost(host string) OptFunc[Server] {
	return func(server *Server) {
		server.host = host
	}
}

func WithPort(port int) OptFunc[Server] {
	return func(server *Server) {
		server.port = port
	}
}

func WithTimeout(timeout time.Duration) OptFunc[Server] {
	return func(server *Server) {
		server.timeout = timeout
	}
}

func WithTLS(server *Server) {
	server.useTLS = true
}

func WithoutTLS(server *Server) {
	server.useTLS = false
}

func NewServer(options ...OptFunc[Server]) *Server {
	// Default server
	server := &Server{
		host:    "127.0.0.1",
		port:    6969,
		timeout: 5 * time.Second,
		useTLS:  false,
	}

	// Decorate with options
	for _, opt := range options {
		opt(server)
	}

	return server
}

func main() {
	server1 := NewServer(
		WithHost("localhost"),
		WithPort(8000),
		WithTLS,
	)

	server2 := NewServer(
		WithTimeout(3*time.Second),
		WithoutTLS,
	)

	fmt.Println(server1)
	fmt.Println(server2)
}
