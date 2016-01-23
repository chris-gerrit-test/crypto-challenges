package main

import (
	"crypto/hmac"
	crand "crypto/rand"
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

func Bob(d dh.DiffieHellman, pub dh.Element) []byte {
	k := d.SharedSecret(pub)
	return Sign(secretMessage, k)
}

var rnd = rand.New(rand.NewSource(82480))

var keySize int = 512

func choosePrime(bits int, m, s *big.Int, avoidFactors map[int64]bool) (*big.Int, map[int64]bool) {
	p := new(big.Int)
	var pFactors map[int64]bool
	for {
		p.SetInt64(2)
		pFactors = make(map[int64]bool)
		pFactors[2] = true
		for i := 0; i < bits/16; i++ {
			f := int64(2)
			for {
				f = rnd.Int63n(1 << 18)
				if ok := avoidFactors[f]; ok {
					continue
				}
				ok := pFactors[f]
				if !ok && big.NewInt(f).ProbablyPrime(32) {
					break
				}
			}
			pFactors[f] = true
			p.Mul(p, big.NewInt(f))
		}
		p.Add(p, big.NewInt(1))
		if p.ProbablyPrime(32) {
			pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
			primRoot := true
			for f, _ := range pFactors {
				pow := new(big.Int).Div(pMinus1, big.NewInt(f))
				check1 := new(big.Int).Exp(m, pow, p)
				check2 := new(big.Int).Exp(s, pow, p)
				if check1.Cmp(big.NewInt(1)) == 0 || check2.Cmp(big.NewInt(1)) == 0 {
					primRoot = false
				}
			}
			if primRoot {
				return p, pFactors
			}
		}
	}
}

func crt(N *big.Int, moduli map[int64]int64) *big.Int {
	// From Wikipedia CRT page
	x := new(big.Int)
	total := N
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
	return x
}

// Find x such that e = g^x (mod N)
func pohlig(e, g, N *big.Int, nFactors map[int64]bool) *big.Int {
	moduli := make(map[int64]int64)
	s := new(big.Int)
	//log.Printf("Pohlig %d", e)
	NMinus1 := new(big.Int).Sub(N, big.NewInt(1))
	for f, _ := range nFactors {
		//log.Printf("%d", f)
		target := new(big.Int).Div(NMinus1, big.NewInt(f))
		base := new(big.Int).Exp(g, target, N)
		target.Exp(e, target, N)
		seen := make(map[int64]bool)
		for m := new(big.Int); ; m.Add(m, big.NewInt(1)) {
			s.Exp(base, m, N)
			seen[s.Int64()] = true
			if s.Cmp(target) == 0 {
				moduli[f] = m.Int64()
				break
			}
			if m.Cmp(new(big.Int).Mul(big.NewInt(f), big.NewInt(2))) >= 0 {
				log.Panicf("Got too high: %d tried", len(seen))
			}
		}
	}
	return crt(NMinus1, moduli)
}

func main() {
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	a := big.NewInt(-95051)
	b := big.NewInt(11279326)
	G := dh.NewEllipticCurve(a, b, p)

	gx := big.NewInt(182)
	gy, _ := new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	g := dh.NewEllipticCurveElement(G, gx, gy)
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	GG := dh.NewGeneratedGroup(G, g, *q)

	log.Printf("========== PART 1 ==========")
	alice := dh.NewECDSA(GG)
	log.Printf("Alice's group identity: %s", alice.Group().Identity())
	log.Printf("Alice's group generator: %s", alice.Group().Generator())
	log.Printf("Alice's generator^q: %s", alice.Group().Pow(alice.Group().Generator(), q))
	r, s := alice.Sign(secretMessage)
	log.Printf("r=%d s=%d", r, s)
	log.Printf("Verifies: %v",
		dh.ECDSAVerify(secretMessage, r, s, alice))
	log.Printf("Verifies with different message: %v",
		dh.ECDSAVerify([]byte("xxx"), r, s, alice))
	r.Add(r, big.NewInt(1))
	r.Mod(r, p)
	log.Printf("Verifies after r += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, alice))
	s.Add(s, big.NewInt(1))
	s.Mod(s, p)
	log.Printf("Verifies after r and s += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, alice))
	r.Sub(r, big.NewInt(1))
	r.Mod(r, p)
	log.Printf("Verifies after s += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, alice))
	s.Sub(s, big.NewInt(1))
	s.Mod(s, p)
	log.Printf("Verifies original: %v",
		dh.ECDSAVerify(secretMessage, r, s, alice))

	eve := dh.FindVerifyingECDSA(secretMessage, r, s, GG)
	log.Printf("Eve's group identity: %s", eve.Group().Identity())
	log.Printf("Eve's group generator: %s", eve.Group().Generator())
	log.Printf("Eve's generator^q: %s", eve.Group().Pow(eve.Group().Generator(), q))
	log.Printf("Eve verifies: %v",
		dh.ECDSAVerify(secretMessage, r, s, eve))

	log.Printf("========== PART 2 ==========")
	pi, err := crand.Prime(crand.Reader, keySize/2)
	if err != nil {
		panic(err)
	}
	qi, err := crand.Prime(crand.Reader, keySize/2)
	log.Printf("Bob's N=%d", new(big.Int).Mul(qi, pi))
	if err != nil {
		panic(err)
	}
	bob := dh.NewRSA(pi, qi)
	h := sha1.New()
	if n, err := h.Write(secretMessage); n != len(secretMessage) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	// Don't feel like implementing padding again
	z := new(big.Int).SetBytes(e) // "pad(m)" in the problem description
	d := bob.Decrypt(z)
	log.Printf("Bob's signature verifies: %v",
		hmac.Equal(e, bob.Encrypt(d).Bytes()))

	padM := new(big.Int).SetBytes(d)
	eveP, evePFactors := choosePrime(keySize/2, z, padM, nil)
	log.Printf("Eve's p=%d", eveP)
	eveQ, eveQFactors := choosePrime(keySize/2, z, padM, evePFactors)
	log.Printf("Eve's q=%d", eveQ)
	log.Printf("Eve's N=%d", new(big.Int).Mul(eveQ, eveP))

	log.Printf("pad(m) mod %d = %d", eveP, new(big.Int).Mod(padM, eveP))
	ep := pohlig(padM, z, eveP, evePFactors)
	log.Printf("%d^%d = %d (mod %d)", z, ep, new(big.Int).Exp(z, ep, eveP), eveP)

	log.Printf("pad(m) mod %d = %d", eveQ, new(big.Int).Mod(padM, eveQ))
	eq := pohlig(padM, z, eveQ, eveQFactors)
	log.Printf("%d^%d = %d (mod %d)", z, eq, new(big.Int).Exp(z, eq, eveQ), eveQ)
}
