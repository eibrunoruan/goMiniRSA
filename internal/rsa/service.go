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
	if bits < 16 {
		return nil, errors.New("bits deve ser >= 16")
	}

	e := big.NewInt(65537)
	pBits := bits / 2
	qBits := bits - pBits

	for {
		p, err := rand.Prime(rand.Reader, pBits)
		if err != nil {
			return nil, err
		}

		q, err := rand.Prime(rand.Reader, qBits)
		if err != nil {
			return nil, err
		}
		if p.Cmp(q) == 0 {
			continue
		}

		n := new(big.Int).Mul(p, q)
		if n.BitLen() != bits {
			continue
		}

		// phi = (p-1)(q-1)
		pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
		qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
		phi := new(big.Int).Mul(pMinus1, qMinus1)

		d := new(big.Int).ModInverse(e, phi)
		if d == nil {
			continue
		}

		return &Keypair{E: e, D: d, N: n}, nil
	}
}

func (s *rsaService) Encrypt(msg *big.Int, e, n *big.Int) *big.Int {
	return new(big.Int).Exp(msg, e, n)
}

func (s *rsaService) Decrypt(cipher *big.Int, d, n *big.Int) *big.Int {
	return new(big.Int).Exp(cipher, d, n)
}
