package main

import (
	"fmt"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
)

// Mail service : Subscriber example

func runMailService() {
	lineMap := make(dict.BoolMap)
	_addAccount := subscribe(AccountBroker, E_ADD_ACCOUNT, 3, lineMap)
	_toggleAccount := subscribe(ToggleBroker, E_TOGGLE_ACCOUNT, 3, lineMap)

	fmt.Println("[Mail] Service started...")
	for {
		select {
		case account, ok := <-_addAccount.Channel:
			runOrClose(mailAddAccount, account, ok, _addAccount, lineMap)
		case params, ok := <-_toggleAccount.Channel:
			runOrClose(mailToggleAccount, params, ok, _toggleAccount, lineMap)
		}
		// Exit if no more active channels
		if check.AllFalse(dict.Values(lineMap)) {
			break
		}
	}
	fmt.Println("[Mail] Service stopped...")
}

func mailAddAccount(account *Account) {
	fmt.Println("[Mail] Add account:", account)
}

func mailToggleAccount(params *ToggleParams) {
	fmt.Println("[Mail] Toggle account:", params)
}
