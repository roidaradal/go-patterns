package main

import (
	"fmt"
	"time"
)

type Account struct {
	code     string
	name     string
	password string
}

func addAccount(account *Account) error {
	rq := NewRequest("AddAccount")
	time.Sleep(2 * time.Second) // artificial delay
	fmt.Println("Add account:", account)

	AccountBroker.Publish(E_ADD_ACCOUNT, account)
	RequestBroker.Publish(E_END_REQUEST, rq.End())

	return nil
}

func editAccount(account *Account) error {
	rq := NewRequest("EditAccount")
	time.Sleep(2 * time.Second) // artificial delay
	fmt.Println("Edit account:", account)

	AccountBroker.Publish(E_EDIT_ACCOUNT, account)
	RequestBroker.Publish(E_END_REQUEST, rq.End())

	return nil
}

func toggleAccount(params *ToggleParams) error {
	rq := NewRequest("ToggleAccount")
	time.Sleep(1 * time.Second) // artificial delay
	fmt.Println("Toggle account:", params)

	ToggleBroker.Publish(E_TOGGLE_ACCOUNT, params)
	RequestBroker.Publish(E_END_REQUEST, rq.End())

	return nil
}
