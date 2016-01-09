package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
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
	b := big.NewInt(11279326)
	G := dh.NewEllipticCurve(a, b, p)

	gx := big.NewInt(182)
	gy, _ := new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	g := dh.NewEllipticCurveElement(G, gx, gy)
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	GG := dh.NewGeneratedGroup(G, g, *q)

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
	log.Printf("Eve verifies: %v",
		dh.ECDSAVerify(secretMessage, r, s, eve))
	log.Printf("Eve's group identity: %s", eve.Group().Identity())
	log.Printf("Eve's group generator: %s", eve.Group().Generator())
	log.Printf("Eve's generator^q: %s", eve.Group().Pow(eve.Group().Generator(), q))

	pi, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	qi, err := rand.Prime(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	bob := dh.NewRSA(pi, qi)
	h := sha1.New()
	if n, err := h.Write(secretMessage); n != len(secretMessage) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	z := new(big.Int).SetBytes(e)
	d := bob.Decrypt(z)
	log.Printf("Bob's signature verifies: %v",
		hmac.Equal(e, bob.Encrypt(d).Bytes()))

	// Find a new
}
