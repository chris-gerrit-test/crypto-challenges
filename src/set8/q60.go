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

func Bob(d dh.DiffieHellman, pub dh.Element) []byte {
	k := d.SharedSecret(pub)
	return Sign(secretMessage, k)
}

var rnd = rand.New(rand.NewSource(3410))

func main() {
	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	n, _ := new(big.Int).SetString("233970423115425145498902418297807005944", 10)
	a := big.NewInt(-95051)
	b := big.NewInt(11279326)
	EC := dh.NewEllipticCurve(a, b, p)

	A := big.NewInt(534)
	B := big.NewInt(1)
	MC := dh.NewMontgomeryCurve(A, B, p)

	mg := dh.NewMontgomeryCurveElement(MC, big.NewInt(4))

	log.Printf("%s^1 = %s", mg, MC.Pow(mg, big.NewInt(1)))
	log.Printf("%s^2 = %s", mg, MC.Pow(mg, big.NewInt(2)))
	log.Printf("%s^%d = %s", mg, q, MC.Pow(mg, q))
	log.Printf("%s^%d = %s", mg, n, MC.Pow(mg, n))

	log.Printf("Generator converted to Weierstrass: %s", dh.Q60MontgomeryToWeierstrass(mg, EC))

	for i := 0; i < 10; i++ {
		w := EC.Random()
		m := dh.Q60WeierstrassToMontgomery(w, MC)
		pow := new(big.Int).Rand(rnd, n)
		w2 := EC.Pow(w, pow)
		m2 := MC.Pow(m, pow)
		log.Printf("%s^%d = %s -> %s", w, pow, w2, dh.Q60WeierstrassToMontgomery(w2, MC))
		log.Printf("%s^%d = %s -> %s", m, pow, m2, dh.Q60MontgomeryToWeierstrass(m2, EC))
	}
}
