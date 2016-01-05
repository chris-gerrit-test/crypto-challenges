package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"log"
	"math/big"
	"math/rand"

	"github.com/hundt/dh"
)

var secretMessage = []byte("Swinging through your town like your neighborhood Spider-Man")

func Sign(msg []byte, key dh.Element) []byte {
	h := hmac.New(sha1.New, []byte(key.String()))
	return h.Sum(msg)
}

func Alice(d dh.DiffieHellman, pub dh.Element) []byte {
	k := d.SharedSecret(pub)
	return Sign(secretMessage, k)
}

var rnd = rand.New(rand.NewSource(3412))

func main() {
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	groupOrder, _ := new(big.Int).SetString("233970423115425145498902418297807005944", 10)
	a := big.NewInt(-95051)
	b := big.NewInt(11279326)
	EC := dh.NewEllipticCurve(a, b, p)
	gx := big.NewInt(182)
	gy, _ := new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	g := dh.NewEllipticCurveElement(EC, gx, gy)

	MC := dh.NewMontgomeryCurve(big.NewInt(534), big.NewInt(1), p)
	mg := dh.NewMontgomeryCurveElement(MC, big.NewInt(4))
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	GG := dh.NewGeneratedGroup(MC, mg, *q)

	twistOrder := new(big.Int).Mul(p, big.NewInt(2))
	twistOrder.Add(twistOrder, big.NewInt(2))
	twistOrder.Sub(twistOrder, groupOrder)

	factors := dh.FindFactors(twistOrder, twistOrder, q, MC)

	total := big.NewInt(1)

	for factor, _ := range factors {
		total.Mul(total, big.NewInt(factor))
	}

	// Find an element that is order=total. We will use this to guess which
	// of the one or two remainders for each small prime is correct.
	log.Printf("Finding an element of order %d...", total)
	pow := new(big.Int)
	pow.Div(twistOrder, total)
	var h dh.Element = nil
	for {
		h = MC.Pow(MC.Random(), pow)
		if h.Cmp(MC.Identity()) != 0 {
			bad := false
			for factor, _ := range factors {
				if MC.Pow(h, big.NewInt(factor)).Cmp(MC.Identity()) == 0 {
					bad = true
					break
				}
			}
			if !bad {
				log.Printf("%s^%d == %s", h, total, MC.Pow(h, total))
				break
			}
		}
	}

	d := dh.NewDiffieHellman(GG)
	mac := Alice(d, h)

	// Determine the possible remainders for each small prime. There are usually
	// two because (u, v) and (u, -v) are conflated by the ladder.
	type factorModuli struct {
		factor         int64
		possibleModuli []int64
	}
	smallFactorModuli := make([]*factorModuli, 0)
	for factor, h := range factors {
		fm := &factorModuli{factor, make([]int64, 0)}
		mac := Alice(d, h)
		log.Printf("Guessing the shared secret in the subgroup of order %d", factor)
		found := false
		for i := int64(0); i < factor; i++ {
			k := MC.Pow(h, big.NewInt(i))
			if hmac.Equal(mac, Sign(secretMessage, k)) {
				found = true
				fm.possibleModuli = append(fm.possibleModuli, i)
				if i != 0 {
					// If h^i = (u, v) then (u, -v) = (h^i)^-1 = h^-i
					fm.possibleModuli = append(fm.possibleModuli, factor-i)
				}
				break
			}
		}
		if !found {
			log.Fatal("Could not guess the shared secret")
		}
		smallFactorModuli = append(smallFactorModuli, fm)
	}

	// Figure out which one of the various choices for each small remainder is correct.
	// For each permutation:
	//  - use CRT to compute remainder x mod "total"
	//  - check if the saved MAC corresponds to h^x.
	indices := make([]int, len(smallFactorModuli))
	bigModuli := make([]*big.Int, 0)
	for {
		carry := true
		for i := 0; carry && i < len(indices); i++ {
			indices[i] = (indices[i] + 1) % len(smallFactorModuli[i].possibleModuli)
			carry = indices[i] == 0
		}
		moduli := make(map[int64]int64)
		for i, fm := range smallFactorModuli {
			moduli[fm.factor] = fm.possibleModuli[indices[i]]
		}
		// From Wikipedia CRT page
		x := new(big.Int)
		for n, a := range moduli {
			N := new(big.Int).Set(total)
			N.Div(N, big.NewInt(n))
			N.ModInverse(N, big.NewInt(n))
			N.Mul(N, total)
			N.Div(N, big.NewInt(n))
			N.Mul(N, big.NewInt(a))
			x.Add(x, N)
			x.Mod(x, total)
		}
		k := MC.Pow(h, x)
		if hmac.Equal(mac, Sign(secretMessage, k)) {
			log.Printf("x mod r may be %s [r = %s]", x, total)
			bigModuli = append(bigModuli, x)
			if x.Cmp(new(big.Int)) != 0 {
				// Again there are probably two answers, because
				// h^x = h^-x using the montgomery ladder.
				y := new(big.Int).Sub(total, x)
				log.Printf("x mod r may be %s [r = %s]", y, total)
				bigModuli = append(bigModuli, y)
			}
			break
		}
	}

	// We now have two possible remainders mod "total". For each of those,
	// use the kangaroo algorithm to try to guess the secret. Note that
	// the public key that Alice gives us will *also* correspond to two
	// different points (u, v) and (u, -v) so we need to try plugging both
	// of them in and see which one works.
	ch := make(chan bool)
	for _, x := range bigModuli {
		for _, negate := range []bool{true, false} {
			go func(x *big.Int, negate bool) {
				n := new(big.Int)
				n.Sub(q, x)

				// Convert to Weierstrass form so Op is implemented and so
				// we can specify which of the two collapsed points (u, v)
				// and (u, -v) to use for this trial.
				w1 := dh.Q60MontgomeryToWeierstrass(d.PublicKey(), EC)
				if negate {
					dh.NegateWeierstrass(w1, EC)
				}
				w2 := EC.Pow(g, n)
				hp := EC.Op(w1, w2)
				gp := EC.Pow(g, total)
				B := new(big.Int)
				B.Div(q, total)

				log.Printf("Searching for index of %s in 0..%d", hp, B.Int64())
				m := dh.Pollard(EC, gp, hp, 0, int(B.Int64()))
				if m == nil {
					log.Printf("The wild kangaroo escaped!")
					return
				} else {
					m.Mod(m, q)
					log.Printf("index = %s", m)
					check := EC.Pow(gp, m)
					log.Printf("%s^%d mod %s = %s", gp, m, p, check)
					if check.Cmp(hp) != 0 {
						log.Printf("ERROR! Expected %s but got %s", hp, check)
					}
				}

				m.Mul(m, total)
				x.Add(x, m)

				log.Printf("Predicted key: %d", x)
				log.Printf("%s", d)
				ch <- true
			}(x, negate)
		}
	}

	<-ch
}
