package rsa

import "math/big"

type Controller struct {
	svc       Service
	presenter Presenter
	attacker  *RSAAttack
	bits      int
}

func NewController(s Service, p Presenter, bits int) *Controller {
	return &Controller{svc: s, presenter: p, attacker: &RSAAttack{}, bits: bits}
}

func (c *Controller) Execute(val uint16) {
	// 1. Setup
	keys, _ := c.svc.GenerateKeys(c.bits)
	c.presenter.ShowKeys(keys)

	// 2. Fluxo Normal
	m := new(big.Int).SetUint64(uint64(val))
	cipher := c.svc.Encrypt(m, keys.E, keys.N)
	c.presenter.ShowProcess("ENCRIPTANDO", m, cipher)

	plain := c.svc.Decrypt(cipher, keys.D, keys.N)
	c.presenter.ShowProcess("DECIFRANDO", cipher, plain)

	// 3. Ataque (Bônus)
	attack, err := c.attacker.Break(keys.E, keys.N)
	if err != nil {
		c.presenter.ShowError("Falha ao fatorar o modulo: " + err.Error())
		return
	}

	hackedPlain := c.svc.Decrypt(cipher, attack.D, keys.N)
	c.presenter.ShowBreakResult(attack.D, hackedPlain, attack.Duration, attack.Backend)
}
