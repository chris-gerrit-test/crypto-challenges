package main

import (
	"crypto/hmac"
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

	ecdsa := dh.NewECDSA(GG)
	log.Printf("%s", ecdsa.PublicKey())
	r, s := ecdsa.Sign(secretMessage)
	log.Printf("r=%d s=%d", r, s)
	log.Printf("Verifies: %v",
		dh.ECDSAVerify(secretMessage, r, s, ecdsa.PublicKey(), GG))
	log.Printf("Verifies with different message: %v",
		dh.ECDSAVerify([]byte("xxx"), r, s, ecdsa.PublicKey(), GG))
	r.Add(r, big.NewInt(1))
	r.Mod(r, p)
	log.Printf("Verifies after r += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, ecdsa.PublicKey(), GG))
	s.Add(s, big.NewInt(1))
	s.Mod(s, p)
	log.Printf("Verifies after r and s += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, ecdsa.PublicKey(), GG))
	r.Sub(r, big.NewInt(1))
	r.Mod(r, p)
	log.Printf("Verifies after s += 1: %v",
		dh.ECDSAVerify(secretMessage, r, s, ecdsa.PublicKey(), GG))
}
