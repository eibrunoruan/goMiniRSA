package rsa

import (
	"context"
	"math/big"
	"runtime"
	"time"
)

type RSAAttack struct{}

func (a *RSAAttack) Break(e, n *big.Int) (*big.Int, time.Duration) {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan *big.Int)
	
	// 1. Pre-Sieve: Checa se é par ou divisível por primos pequenos
	smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23}
	for _, p := range smallPrimes {
		bp := big.NewInt(p)
		if new(big.Int).Mod(n, bp).Sign() == 0 {
			return a.deriveD(e, bp, n), time.Since(start)
		}
	}

	// 2. Paralelismo Massivo: Dispara um caçador por core da CPU
	numCPUs := runtime.NumCPU()
	for i := 0; i < numCPUs; i++ {
		seed := int64(i + 1)
		go func(s int64) {
			if p := pollardBrent(ctx, n, s); p != nil {
				select {
				case resultChan <- p:
				case <-ctx.Done():
				}
			}
		}(seed)
	}

	p := <-resultChan
	cancel()

	return a.deriveD(e, p, n), time.Since(start)
}

func (a *RSAAttack) deriveD(e, p, n *big.Int) *big.Int {
	q := new(big.Int).Div(n, p)
	phi := new(big.Int).Mul(
		new(big.Int).Sub(p, big.NewInt(1)),
		new(big.Int).Sub(q, big.NewInt(1)),
	)
	return new(big.Int).ModInverse(e, phi)
}

// Algoritmo de Brent: Versão otimizada do Pollard's Rho
func pollardBrent(ctx context.Context, n *big.Int, seed int64) *big.Int {
	y, c, m := big.NewInt(seed), big.NewInt(seed), big.NewInt(seed)
	g, r, q := big.NewInt(1), big.NewInt(1), big.NewInt(1)
	x, ys := big.NewInt(0), big.NewInt(0)
	one := big.NewInt(1)

	for g.Cmp(one) == 0 {
		x.Set(y)
		for i := big.NewInt(0); i.Cmp(r) < 0; i.Add(i, one) {
			y.Mod(new(big.Int).Add(new(big.Int).Mul(y, y), c), n)
		}
		
		k := big.NewInt(0)
		for k.Cmp(r) < 0 && g.Cmp(one) == 0 {
			select {
			case <-ctx.Done(): return nil
			default:
			}
			
			ys.Set(y)
			limit := new(big.Int)
			if r.Cmp(new(big.Int).Sub(r, k)) < 0 { limit.Set(r) } else { limit.Sub(r, k) }
			if m.Cmp(limit) < 0 { limit.Set(m) }

			for i := big.NewInt(0); i.Cmp(limit) < 0; i.Add(i, one) {
				y.Mod(new(big.Int).Add(new(big.Int).Mul(y, y), c), n)
				q.Mod(q.Mul(q, new(big.Int).Abs(new(big.Int).Sub(x, y))), n)
			}
			g.GCD(nil, nil, q, n)
			k.Add(k, m)
		}
		r.Mul(r, big.NewInt(2))
	}

	if g.Cmp(n) == 0 {
		for {
			ys.Mod(new(big.Int).Add(new(big.Int).Mul(ys, ys), c), n)
			g.GCD(nil, nil, new(big.Int).Abs(new(big.Int).Sub(x, ys)), n)
			if g.Cmp(one) > 0 { break }
		}
	}
	return g
}