package main

import (
	"mini-rsa/internal/rsa"
)

func main() {
	service := rsa.NewService()
	presenter := rsa.NewPresenter()
	
	app := rsa.NewController(service, presenter)

	// Valor de 16 bits para teste
	app.Execute(54321)
}