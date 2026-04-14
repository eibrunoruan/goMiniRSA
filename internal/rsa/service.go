package rsa

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Keypair armazena as chaves geradas
type Keypair struct {
	E, D, N *big.Int
}

// Service define a interface de negócio (Model)
type Service interface {
	GenerateKeys(bits int) (*Keypair, error)
	Encrypt(msg *big.Int, pubE, n *big.Int) *big.Int
	Decrypt(cipher *big.Int, privD, n *big.Int) *big.Int
}

type rsaService struct{}

func NewService() Service {
	return &rsaService{}
}

func (s *rsaService) GenerateKeys(bits int) (*Keypair, error) {
	// Para suportar 16 bits, geramos primos de (bits/2)+1
	p, err := rand.Prime(rand.Reader, (bits/2)+1)
	if err != nil {
		return nil, err
	}
	q, err := rand.Prime(rand.Reader, (bits/2)+1)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).Mul(p, q)
	
	// phi = (p-1)(q-1)
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)

	// e = 65537 é o padrão de mercado (F4 de Fermat)
	e := big.NewInt(65537)
	
	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		return nil, errors.New("erro ao calcular inverso modular")
	}

	return &Keypair{E: e, D: d, N: n}, nil
}

func (s *rsaService) Encrypt(msg *big.Int, e, n *big.Int) *big.Int {
	return new(big.Int).Exp(msg, e, n)
}

func (s *rsaService) Decrypt(cipher *big.Int, d, n *big.Int) *big.Int {
	return new(big.Int).Exp(cipher, d, n)
}