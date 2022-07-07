package main

import (
	customerOrderApp "br.com.thyagoluciano.poc/app/customer-order-app/cmd"
	paymentApp "br.com.thyagoluciano.poc/app/payment-app/cmd"
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

	go func() {
		paymentApp.Start()
	}()

	wg.Wait()
}
