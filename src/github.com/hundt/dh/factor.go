package dh

import (
	"log"
	"math"
	"math/big"
)

func FindFactors(j, groupSize, target *big.Int, G Group) map[int64]Element {
	zero := new(big.Int)
	factors := make(map[int64]Element, 0)
	total := big.NewInt(1)
	for prime := int64(2); total.Cmp(target) < 0; prime++ {
		// Some quick exits that cover the majority of cases
		if prime > 11 {
			if prime%2 == 0 || prime%3 == 0 || prime%5 == 0 || prime%7 == 0 || prime%11 == 0 {
				continue
			}
		}
		pr := big.NewInt(prime)
		pr.Rem(j, pr)
		if pr.Cmp(zero) == 0 {
			if !big.NewInt(prime).ProbablyPrime(20) {
				continue
			}
			j2 := new(big.Int).Set(j)
			j2.Div(j, big.NewInt(prime))
			pr.Rem(j2, big.NewInt(prime))
			if pr.Cmp(zero) == 0 {
				log.Printf("Skipping double divisor %d", prime)
			} else {
				log.Printf("Found divisor %d", prime)
				factors[prime] = nil
				total.Mul(total, big.NewInt(prime))
			}
		}
		if prime > 1<<22 {
			log.Printf("Giving up with total=%s", total)
			break
		}
	}

	for factor, _ := range factors {
		var h Element = nil
		groupSize := new(big.Int).Set(groupSize)
		pow := new(big.Int)
		pow.Div(groupSize, big.NewInt(factor))
		log.Printf("Finding an element of order %d...", factor)
		for {
			h = G.Pow(G.Random(), pow)
			if h.Cmp(G.Identity()) != 0 {
				//log.Printf("%s^%d == %s", h, factor, G.Pow(h, big.NewInt(factor)))
				break
			}
		}
		factors[factor] = h
	}

	return factors
}

func FindCoFactors(q, n *big.Int, G Group) map[int64]Element {
	j := new(big.Int)
	j.Sub(n, big.NewInt(1))
	j.Div(j, q)

	groupSize := new(big.Int).Sub(n, big.NewInt(1))

	return FindFactors(j, groupSize, q, G)
}

func JumpMean(G Group, k int) float64 {
	trials := 1000
	total := 0.0
	for i := 0; i < trials; i++ {
		total += float64(G.Random().Jump(k))
	}
	return total / float64(trials)
}

func Pollard(G Group, g, h Element, a, b int) *big.Int {
	k := int(math.Log2(float64(b-a)))/2 + 2
	N := int(4 * JumpMean(G, k))
	log.Printf("k = %d, N = %d", k, N)
	tameDistance := 0
	tamePos := G.Pow(g, big.NewInt(int64(b)))
	log.Printf("Setting the trap")
	for i := 0; i < N; i++ {
		jump := tamePos.Jump(k)
		tameDistance += jump
		tamePos = G.Op(tamePos, G.Pow(g, big.NewInt(int64(jump))))
	}

	wildDistance := 0
	wildPos := h
	threshold := b - a + tameDistance
	log.Printf("Waiting to catch the wild kangaroo")
	nextProgress := threshold / 10
	for wildDistance < threshold {
		jump := wildPos.Jump(k)
		wildDistance += jump
		wildPos = G.Op(wildPos, G.Pow(g, big.NewInt(int64(jump))))
		if tamePos.Cmp(wildPos) == 0 {
			return big.NewInt(int64(b + tameDistance - wildDistance))
		}
		if wildDistance > nextProgress {
			nextProgress += threshold / 10
		}
	}
	return nil
}
