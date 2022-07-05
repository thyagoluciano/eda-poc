package main

import (
	customerOrderApp "br.com.thyagoluciano.poc/app/customer-order-app/cmd"
	shoppingCartApp "br.com.thyagoluciano.poc/app/shopping-cart-app/cmd"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)

	wg.Add(10)

	go func() {
		shoppingCartApp.Start()
	}()

	go func() {
		customerOrderApp.Start()
	}()

	wg.Wait()
}
