package main

import (
	"fmt"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
)

// Cache service : Subscriber example

func runCacheService() {
	lineMap := make(dict.BoolMap)
	_addAccount := subscribe(AccountBroker, E_ADD_ACCOUNT, 5, lineMap)
	_editAccount := subscribe(AccountBroker, E_EDIT_ACCOUNT, 5, lineMap)
	_toggleAccount := subscribe(ToggleBroker, E_TOGGLE_ACCOUNT, 2, lineMap)

	fmt.Println("[Cache] Service started...")
	for {
		select {
		case account, ok := <-_addAccount.Channel:
			runOrClose(cacheAddAccount, account, ok, _addAccount, lineMap)
		case account, ok := <-_editAccount.Channel:
			runOrClose(cacheEditAccount, account, ok, _editAccount, lineMap)
		case params, ok := <-_toggleAccount.Channel:
			runOrClose(cacheToggleAccount, params, ok, _toggleAccount, lineMap)
		}
		// Exit if no more active channels
		if check.AllFalse(dict.Values(lineMap)) {
			break
		}
	}
	fmt.Println("[Cache] Service stopped...")
}

func cacheAddAccount(account *Account) {
	fmt.Println("[Cache] Add account:", account)
}

func cacheEditAccount(account *Account) {
	fmt.Println("[Cache] Edit account:", account)
}

func cacheToggleAccount(params *ToggleParams) {
	fmt.Println("[Cache] Toggle account:", params)
}
