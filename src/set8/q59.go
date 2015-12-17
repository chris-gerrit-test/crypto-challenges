package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"
	"math/big"

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

func main() {
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	a := big.NewInt(-95051)
	//a.Sub(p, a)
	b := big.NewInt(11279326)
	G := dh.NewEllipticCurve(a, b, p)

	gx := big.NewInt(182)
	gy, _ := new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	g := dh.NewEllipticCurveElement(G, gx, gy)
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	GG := dh.NewGeneratedGroup(G, g, *q)

	d1 := dh.NewDiffieHellman(GG)
	d2 := dh.NewDiffieHellman(GG)

	fmt.Printf("=== PHASE I: TEST DIFFIE-HELLMAN WORKS ===\n")

	log.Printf("Alice's DH: %s", d1)
	log.Printf("Bob's DH: %s", d2)

	log.Printf("Bob's shared secret from Alice's public key: %s", d2.SharedSecret(d1.PublicKey()))
	log.Printf("Alice's shared secret from Bob's public key: %s", d1.SharedSecret(d2.PublicKey()))

	fmt.Printf("=== PHASE II: BREAK DIFFIE-HELLMAN ===\n")

	weakSubgroups := []dh.Group{
		dh.NewEllipticCurve(big.NewInt(-95051), big.NewInt(210), p),
		dh.NewEllipticCurve(big.NewInt(-95051), big.NewInt(504), p),
		dh.NewEllipticCurve(big.NewInt(-95051), big.NewInt(727), p),
	}
	weakSubgroupOrders := make([]*big.Int, 3)
	weakSubgroupOrders[0], _ = new(big.Int).SetString("233970423115425145550826547352470124412", 10)
	weakSubgroupOrders[1], _ = new(big.Int).SetString("233970423115425145544350131142039591210", 10)
	weakSubgroupOrders[2], _ = new(big.Int).SetString("233970423115425145545378039958152057148", 10)

	allFactors := make(map[int64]dh.Element)

	for i, SG := range weakSubgroups {
		order := weakSubgroupOrders[i]
		factors := dh.FindFactors(order, order, q, SG)
		for factor, h := range factors {
			if _, ok := allFactors[factor]; !ok {
				allFactors[factor] = h
			}
		}
	}

	moduli := make(map[int64]int64)
	total := big.NewInt(1)

	d := dh.NewDiffieHellman(GG)
	for factor, h := range allFactors {
		mac := Bob(d, h)
		log.Printf("Guessing the shared secret in the subgroup of order %d", factor)
		found := false
		for i := int64(1); i <= factor; i++ {
			k := G.Pow(h, big.NewInt(i))
			if hmac.Equal(mac, Sign(secretMessage, k)) {
				found = true
				moduli[factor] = i
				break
			}
		}
		if found {
			total.Mul(total, big.NewInt(factor))
			if total.Cmp(q) > 0 {
				log.Printf("Found enough moduli")
				break
			}
		} else {
			// This can happen if a prime factor of one of the weak subgroups' order
			// is a double divisor in another one's order.
			log.Printf("Could not guess the shared secret!")
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
	log.Printf("Predicted key: %d", x)

	log.Printf("%s", d)
}
