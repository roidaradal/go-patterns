package main

import "fmt"

func runCacheService() {
	addAccountCh := AccountBroker.Subscribe(E_ADD_ACCOUNT, 5)
	editAccountCh := AccountBroker.Subscribe(E_EDIT_ACCOUNT, 5)
	toggleAccountCh := ToggleBroker.Subscribe(E_TOGGLE_ACCOUNT, 2)

	fmt.Println("[Cache] Service started...")
	for {
		select {
		case account := <-addAccountCh:
			cacheAddAccount(account)
		case account := <-editAccountCh:
			cacheEditAccount(account)
		case params := <-toggleAccountCh:
			cacheToggleAccount(params)
		}
	}
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
