package rsa

import "math/big"

type Controller struct {
	svc       Service
	presenter Presenter
	attacker  *RSAAttack
}

func NewController(s Service, p Presenter) *Controller {
	return &Controller{svc: s, presenter: p, attacker: &RSAAttack{}}
}

func (c *Controller) Execute(val uint16) {
	// 1. Setup
	keys, _ := c.svc.GenerateKeys(16)
	c.presenter.ShowKeys(keys)

	// 2. Fluxo Normal
	m := new(big.Int).SetUint64(uint64(val))
	cipher := c.svc.Encrypt(m, keys.E, keys.N)
	c.presenter.ShowProcess("ENCRIPTANDO", m, cipher)

	plain := c.svc.Decrypt(cipher, keys.D, keys.N)
	c.presenter.ShowProcess("DECIFRANDO", cipher, plain)

	// 3. Ataque (Bônus)
	stolenD, duration := c.attacker.Break(keys.E, keys.N)
	hackedPlain := c.svc.Decrypt(cipher, stolenD, keys.N)
	c.presenter.ShowBreakResult(stolenD, hackedPlain, duration)
}