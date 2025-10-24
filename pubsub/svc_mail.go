package main

import "fmt"

func runMailService() {
	addAccountCh := AccountBroker.Subscribe(E_ADD_ACCOUNT, 3)
	toggleAccountCh := ToggleBroker.Subscribe(E_TOGGLE_ACCOUNT, 3)

	fmt.Println("[Mail] Service started...")
	for {
		select {
		case account := <-addAccountCh:
			mailAddAccount(account)
		case params := <-toggleAccountCh:
			mailToggleAccount(params)
		}
	}
}

func mailAddAccount(account *Account) {
	fmt.Println("[Mail] Add account:", account)
}

func mailToggleAccount(params *ToggleParams) {
	fmt.Println("[Mail] Toggle account:", params)
}
