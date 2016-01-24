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

var keySize int = 1024

func choosePrime(bits int, m, s *big.Int, avoidFactors map[int64]bool) (*big.Int, map[int64]bool) {
	p := new(big.Int)
	var pFactors map[int64]bool
	for {
		p.SetInt64(2)
		pFactors = make(map[int64]bool)
		pFactors[2] = true
		for i := 0; i < bits/14; i++ {
			f := int64(2)
			for {
				f = rnd.Int63n(1 << 16)
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
	NMinus1 := new(big.Int).Sub(N, big.NewInt(1))
	for f, _ := range nFactors {
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

// Returns (N, e, d) of the given size that encrypts decBytes to encBytes
func chosenEncrytion(decBytes, encBytes []byte, keySize int) (N, e, d *big.Int) {
	padM := new(big.Int).SetBytes(encBytes)
	sig := new(big.Int).SetBytes(decBytes)
	eveP, evePFactors := choosePrime(keySize/2, padM, sig, nil)
	eveQ, eveQFactors := choosePrime(keySize/2, padM, sig, evePFactors)
	eveN := new(big.Int).Mul(eveP, eveQ)

	ep := pohlig(padM, sig, eveP, evePFactors)

	eq := pohlig(padM, sig, eveQ, eveQFactors)

	evePhi := new(big.Int).Set(eveN)
	evePhi.Sub(evePhi, eveP)
	evePhi.Sub(evePhi, eveQ)
	evePhi.Add(evePhi, big.NewInt(1))

	// From Wikipedia CRT page
	eveModuli := make(map[*big.Int]*big.Int)
	newP := new(big.Int).Sub(eveP, big.NewInt(1))
	newQ := new(big.Int).Sub(eveQ, big.NewInt(1))
	// Remove "2" factor so CRT works.
	newP.Div(newP, big.NewInt(2))
	newQ.Div(newQ, big.NewInt(2))
	eveModuli[newP] = new(big.Int).Mod(ep, newP)
	eveModuli[newQ] = new(big.Int).Mod(eq, newQ)
	if new(big.Int).Mod(eq, big.NewInt(2)).Cmp(new(big.Int).Mod(ep, big.NewInt(2))) != 0 {
		log.Fatal("ep and eq are not compatible and CRT won't work")
	}
	eveModuli[big.NewInt(4)] = new(big.Int).Mod(eq, big.NewInt(4))
	eveE := new(big.Int)
	for n, a := range eveModuli {
		N := new(big.Int).Set(evePhi)
		scratch := new(big.Int)
		scratch.Mod(N, n)
		if scratch.Cmp(new(big.Int)) != 0 {
			log.Fatalf("hmm %d / %d", N, n)
		}
		N.Div(N, n)
		scratch.GCD(nil, nil, N, n)
		if scratch.Cmp(big.NewInt(1)) != 0 {
			log.Fatalf("GCD: %d", scratch)
		}
		N.ModInverse(N, n)
		N.Mul(N, evePhi)
		N.Div(N, n)
		N.Mul(N, a)
		eveE.Add(eveE, N)
		eveE.Mod(eveE, evePhi)
	}

	return eveN, eveE, new(big.Int).ModInverse(eveE, evePhi)
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
	if err != nil {
		panic(err)
	}
	bob := dh.NewRSA(pi, qi)
	log.Printf("Bob's (N, e) = (%d, %d)", new(big.Int).Mul(qi, pi), bob.PublicKey())
	h := sha1.New()
	if n, err := h.Write(secretMessage); n != len(secretMessage) || err != nil {
		log.Fatal("Error calculating hash")
	}
	hashBytes := h.Sum(nil)
	// Don't feel like implementing padding again
	sigBytes := bob.Decrypt(new(big.Int).SetBytes(hashBytes))
	log.Printf("Bob's signature verifies: %v",
		hmac.Equal(hashBytes, bob.Encrypt(sigBytes).Bytes()))

	eveN, eveE, _ := chosenEncrytion(sigBytes, hashBytes, keySize)
	log.Printf("Eve's (N', e') = (%d, %d)", eveN, eveE)
	log.Printf("Eve's signature verifies: %v",
		hmac.Equal(hashBytes, new(big.Int).Exp(new(big.Int).SetBytes(sigBytes), eveE, eveN).Bytes()))

	log.Printf("========== PART 3 ==========")
	original := []byte("Transfer $100 from Chris to the Treasury")
	desired := []byte("Transfer $1,000,000 from the Treasury to Chris")
	pi, err = crand.Prime(crand.Reader, keySize/2)
	if err != nil {
		panic(err)
	}
	qi, err = crand.Prime(crand.Reader, keySize/2)
	if err != nil {
		panic(err)
	}
	bob = dh.NewRSA(pi, qi)
	log.Printf("Bob's (N, e) = (%d, %d)", new(big.Int).Mul(qi, pi), bob.PublicKey())
	enc := bob.Encrypt(original)
	eveN, eveE, eveD := chosenEncrytion(desired, enc.Bytes(), keySize)
	log.Printf("Eve's (N', e', d') = (%d, %d, %d)", eveN, eveE, eveD)
	if eveN.Cmp(new(big.Int).SetBytes(desired)) < 0 {
		log.Panic("Eve's key is too small")
	}
	log.Printf("Decrypted: %s", new(big.Int).Exp(enc, eveD, eveN).Bytes())
}
