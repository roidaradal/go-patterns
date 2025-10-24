package main

import (
	"time"

	"github.com/roidaradal/krap"
)

const (
	E_ADD_ACCOUNT    string = "add-account"
	E_TOGGLE_ACCOUNT string = "toggle-account"
	E_EDIT_ACCOUNT   string = "edit-account"
	E_END_REQUEST    string = "end-request"
)

var (
	AccountBroker = NewPubSub[*Account]()
	RequestBroker = NewPubSub[*Request]()
	ToggleBroker  = NewPubSub[*ToggleParams]()
)

func main() {
	go runCacheService()
	go runMailService()
	go runLogService()

	// Simulate requests
	go func() {
		addAccount(&Account{"roi", "Roi", "abc123"})
		addAccount(&Account{"john", "John", "def456"})
		time.Sleep(5 * time.Second)
		editAccount(&Account{"roi", "Roy", "def666"})
		time.Sleep(3 * time.Second)
		toggleAccount(&ToggleParams{"john", false})
		time.Sleep(5 * time.Second)
		toggleAccount(&ToggleParams{"john", true})
	}()

	cfg := &krap.WebConfig{
		Base: "/api/v1",
		Port: 6666,
	}
	server, address := krap.WebServer(cfg, "dev")
	server.Run(address)
}
