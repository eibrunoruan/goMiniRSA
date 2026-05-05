package main

import (
	"mini-rsa/internal/rsa"
	"os"
	"strconv"
)

func main() {
	service := rsa.NewService()
	presenter := rsa.NewPresenter()
	bits := 16
	if raw := os.Getenv("RSA_BITS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed >= 16 {
			bits = parsed
		}
	}

	app := rsa.NewController(service, presenter, bits)

	// Valor de 16 bits para teste
	app.Execute(54321)
}
