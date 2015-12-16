package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"
	"math/big"

	"github.com/hundt/dh"
)

var secretMessage = []byte("crazy flamboyant for the rap enjoyment")

func Sign(msg []byte, key dh.Element) []byte {
	h := hmac.New(sha1.New, []byte(key.String()))
	return h.Sum(msg)
}

func Bob(d dh.DiffieHellman, pub dh.Element) []byte {
	k := d.SharedSecret(pub)
	return Sign(secretMessage, k)
}

func main() {
	fmt.Printf("=== PHASE I: TEST ALGORITHM ===\n")
	pi := new(big.Int)
	pi.SetString("11470374874925275658116663507232161402086650258453896274534991676898999262641581519101074740642369848233294239851519212341844337347119899874391456329785623", 10)
	gi := new(big.Int)
	gi.SetString("622952335333961296978159266084741085889881358738459939978290179936063635566740258555167783009058567397963466103140082647486611657350811560630587013183357", 10)
	q := new(big.Int)
	q.SetString("335062023296420808191071248367701059461", 10)
	G := dh.NewFiniteGroup(*pi)
	g := dh.NewFiniteElement(G, *gi)
	GG := dh.NewGeneratedGroup(G, g, *q)

	targets := []string{
		"7760073848032689505395005705677365876654629189298052775754597607446617558600394076764814236081991643094239886772481052254010323780165093955236429914607119",
		"9388897478013399550694114614498790691034187453089355259602614074132918843899833277397448144245883225611726912025846772975325932794909655215329941809013733",
	}

	hi := new(big.Int)
	a := 0
	upperBounds := []int{1 << 20, 1 << 40}

	for i, hs := range targets {
		hi.SetString(hs, 10)
		b := upperBounds[i]

		h := dh.NewFiniteElement(G, *hi)

		log.Printf("Searching for order of %s", h)
		x := dh.Pollard(GG, g, h, a, b)
		if x == nil {
			log.Printf("The wild kangaroo escaped!")
		} else {
			x.Mod(x, q)
			log.Printf("order = %s", x)
			check := GG.Pow(g, x)
			log.Printf("%s^%d mod %s == %s", g, x, pi, GG.Pow(g, x))
			if check.Cmp(h) != 0 {
				log.Printf("ERROR! Expected %s but got %s", h, check)
			}
		}
	}

	fmt.Printf("\n=== PHASE II: APPLY WITH COFACTOR ATTACK ===\n")
	factors := dh.FindCoFactors(q, pi, G)

	moduli := make(map[int64]int64)
	total := big.NewInt(1)

	d := dh.NewDiffieHellman(GG)
	for factor, h := range factors {
		total.Mul(total, big.NewInt(factor))
		mac := Bob(d, h)
		log.Printf("Guessing the shared secret in the subgroup of order %d", factor)
		found := false
		for i := int64(1); i <= factor; i++ {
			k := G.Pow(h, big.NewInt(i))
			if hmac.Equal(mac, Sign(secretMessage, k)) {
				//log.Printf("%d^%d", elt, i)
				found = true
				moduli[factor] = i
				break
			}
		}
		if !found {
			panic("Could not guess the shared secret")
		}
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
	log.Printf("x mod r = %s [r = %s]", x, total)

	n := new(big.Int)
	n.Sub(q, x)

	hp := GG.Op(d.PublicKey(), GG.Pow(g, n))
	gp := GG.Pow(g, total)
	B := new(big.Int)
	B.Div(q, total)

	log.Printf("Searching for order of %s", hp)
	m := dh.Pollard(GG, gp, hp, 0, int(B.Int64()))
	if m == nil {
		log.Printf("The wild kangaroo escaped!")
	} else {
		m.Mod(m, q)
		log.Printf("order = %s", m)
		check := GG.Pow(gp, m)
		log.Printf("%s^%d mod %s == %s", gp, m, pi, GG.Pow(gp, m))
		if check.Cmp(hp) != 0 {
			log.Printf("ERROR! Expected %s but got %s", hp, check)
		}
	}

	m.Mul(m, total)
	x.Add(x, m)

	log.Printf("Predicted key: %d", x)
	log.Printf("%s", d)
}
