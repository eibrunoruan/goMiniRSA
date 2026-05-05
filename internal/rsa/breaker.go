package rsa

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type RSAAttack struct{}

type AttackResult struct {
	D        *big.Int
	Factor   *big.Int
	Backend  string
	Duration time.Duration
}

func (a *RSAAttack) Break(e, n *big.Int) (*AttackResult, error) {
	start := time.Now()

	if n.Bit(0) == 0 {
		p := big.NewInt(2)
		return &AttackResult{
			D:        a.deriveD(e, p, n),
			Factor:   p,
			Backend:  "trial-division",
			Duration: time.Since(start),
		}, nil
	}

	// Um crivo simples evita iniciar Pollard-Brent para casos triviais.
	smallPrimes := []int64{
		3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
		53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
	for _, p := range smallPrimes {
		bp := big.NewInt(p)
		if new(big.Int).Mod(n, bp).Sign() == 0 {
			return &AttackResult{
				D:        a.deriveD(e, bp, n),
				Factor:   bp,
				Backend:  "trial-division",
				Duration: time.Since(start),
			}, nil
		}
	}

	if p, backend, err := factorWithMSieve(n); err == nil {
		return &AttackResult{
			D:        a.deriveD(e, p, n),
			Factor:   p,
			Backend:  backend,
			Duration: time.Since(start),
		}, nil
	}

	p := breakWithPollardBrent(n)
	if p == nil {
		return nil, errors.New("nao foi possivel fatorar n")
	}

	return &AttackResult{
		D:        a.deriveD(e, p, n),
		Factor:   p,
		Backend:  "pollard-brent-fallback",
		Duration: time.Since(start),
	}, nil
}

func breakWithPollardBrent(n *big.Int) *big.Int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan *big.Int, 1)

	// Workers independentes reiniciam a caminhada com novas seeds ate achar um fator.
	workerCount := runtime.NumCPU()
	if workerCount < 1 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		go func() {
			for ctx.Err() == nil {
				ySeed, cSeed := randomSeedPair(n)
				if p := pollardBrent(ctx, n, ySeed, cSeed); p != nil {
					select {
					case resultChan <- p:
					case <-ctx.Done():
					}
					return
				}
			}
		}()
	}

	p := <-resultChan
	cancel()
	return p
}

func (a *RSAAttack) deriveD(e, p, n *big.Int) *big.Int {
	q := new(big.Int).Div(n, p)
	phi := new(big.Int).Mul(
		new(big.Int).Sub(p, big.NewInt(1)),
		new(big.Int).Sub(q, big.NewInt(1)),
	)
	return new(big.Int).ModInverse(e, phi)
}

func pollardBrent(ctx context.Context, n, y, c *big.Int) *big.Int {
	one := big.NewInt(1)
	two := big.NewInt(2)

	if new(big.Int).GCD(nil, nil, c, n).Cmp(one) != 0 {
		return nil
	}

	const blockSize int64 = 128

	g := new(big.Int).Set(one)
	q := new(big.Int).Set(one)
	x := new(big.Int)
	ys := new(big.Int)
	diff := new(big.Int)
	tmp := new(big.Int)

	r := int64(1)
	for g.Cmp(one) == 0 {
		x.Set(y)
		for i := int64(0); i < r; i++ {
			if !advanceSequence(ctx, y, c, n, tmp) {
				return nil
			}
		}

		for k := int64(0); k < r && g.Cmp(one) == 0; k += blockSize {
			if ctx.Err() != nil {
				return nil
			}

			ys.Set(y)
			limit := minInt64(blockSize, r-k)
			q.Set(one)

			for i := int64(0); i < limit; i++ {
				if !advanceSequence(ctx, y, c, n, tmp) {
					return nil
				}

				diff.Sub(x, y)
				diff.Abs(diff)
				if diff.Sign() == 0 {
					continue
				}

				q.Mul(q, diff)
				q.Mod(q, n)
			}

			g.GCD(nil, nil, q, n)
		}

		r *= 2
	}

	if g.Cmp(n) == 0 {
		for {
			if !advanceSequence(ctx, ys, c, n, tmp) {
				return nil
			}

			diff.Sub(x, ys)
			diff.Abs(diff)
			g.GCD(nil, nil, diff, n)

			switch g.Cmp(one) {
			case 0:
				continue
			case 1:
				if g.Cmp(n) == 0 {
					return nil
				}
				return g
			}
		}
	}

	if g.Cmp(two) >= 0 && g.Cmp(n) < 0 {
		return g
	}

	return nil
}

func advanceSequence(ctx context.Context, x, c, n, scratch *big.Int) bool {
	if ctx.Err() != nil {
		return false
	}

	scratch.Mul(x, x)
	scratch.Add(scratch, c)
	x.Mod(scratch, n)
	return true
}

func randomSeedPair(n *big.Int) (*big.Int, *big.Int) {
	max := new(big.Int).Sub(n, big.NewInt(3))
	if max.Sign() <= 0 {
		return big.NewInt(2), big.NewInt(1)
	}

	y, err := rand.Int(rand.Reader, max)
	if err != nil {
		return big.NewInt(2), big.NewInt(1)
	}
	y.Add(y, big.NewInt(2))

	c, err := rand.Int(rand.Reader, max)
	if err != nil {
		return y, big.NewInt(1)
	}
	c.Add(c, big.NewInt(1))

	return y, c
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func factorWithMSieve(n *big.Int) (*big.Int, string, error) {
	msievePath, err := exec.LookPath("msieve")
	if err != nil {
		return nil, "", err
	}

	threads := runtime.NumCPU()
	if threads < 1 {
		threads = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	args := []string{
		"-s", filepath.Join(os.TempDir(), "mini-rsa-msieve.dat"),
		"-l", filepath.Join(os.TempDir(), "mini-rsa-msieve.log"),
		"-q",
		"-e",
		"-t", strconv.Itoa(threads),
		n.String(),
	}

	cmd := exec.CommandContext(ctx, msievePath, args...)
	defer os.Remove(filepath.Join(os.TempDir(), "mini-rsa-msieve.dat"))
	defer os.Remove(filepath.Join(os.TempDir(), "mini-rsa-msieve.log"))
	out, err := cmd.Output()
	if err != nil {
		return nil, "", err
	}

	factor, err := parseMSieveFactor(out, n)
	if err != nil {
		return nil, "", err
	}

	return factor, "msieve-ecm", nil
}

func parseMSieveFactor(out []byte, n *big.Int) (*big.Int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		colon := strings.IndexByte(line, ':')
		if colon == -1 {
			continue
		}

		label := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		if !strings.HasPrefix(label, "p") && !strings.HasPrefix(label, "prp") && !strings.HasPrefix(label, "c") {
			continue
		}

		factor, ok := new(big.Int).SetString(value, 10)
		if !ok || factor.Sign() <= 0 {
			continue
		}
		if factor.Cmp(big.NewInt(1)) == 0 || factor.Cmp(n) == 0 {
			continue
		}
		if new(big.Int).Mod(n, factor).Sign() != 0 {
			continue
		}
		return factor, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("msieve nao retornou fator utilizavel")
}
