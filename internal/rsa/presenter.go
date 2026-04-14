package rsa

import (
	"fmt"
	"math/big"
	"time"
)

type Presenter interface {
	ShowKeys(keys *Keypair)
	ShowProcess(label string, input, output *big.Int)
	ShowBreakResult(stolenD, plain *big.Int, duration time.Duration)
}

type ConsolePresenter struct{}

func NewPresenter() Presenter {
	return &ConsolePresenter{}
}

func (p *ConsolePresenter) ShowKeys(keys *Keypair) {
	fmt.Println("\n--- [PASSO 1] GERAÇÃO DE CHAVES ---")
	fmt.Printf("Pública (e): %s\n", keys.E)
	fmt.Printf("Privada (d): %s\n", keys.D)
	fmt.Printf("Módulo  (n): %s\n", keys.N)
	fmt.Println("----------------------------------")
}

func (p *ConsolePresenter) ShowProcess(label string, input, output *big.Int) {
	fmt.Printf("%-15s | In: %-6s | Out: %s\n", label, input, output)
}

func (p *ConsolePresenter) ShowBreakResult(stolenD, plain *big.Int, duration time.Duration) {
	fmt.Println("\n--- [PASSO 2] BÔNUS: QUEBRA POR FATORAÇÃO ---")
	fmt.Printf("Chave 'd' Descoberta: %s\n", stolenD)
	fmt.Printf("Mensagem Decifrada:   %s\n", plain)
	fmt.Printf("Tempo de Execução:    %v\n", duration)
	fmt.Println("------------------------------------------")
}