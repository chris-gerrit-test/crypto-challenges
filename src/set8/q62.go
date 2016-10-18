package main

import (
	"crypto/sha1"
	"log"
	"math/big"
	"os"
	"runtime/pprof"

	"github.com/hundt/dh"
)

func transform(msg []byte, r, s, q *big.Int, bias uint) (*big.Int, *big.Int) {
	h := sha1.New()
	if n, err := h.Write(msg); n != len(msg) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	z := new(big.Int).SetBytes(e)
	z.Mod(z, q)

	t := big.NewInt(1 << bias)
	t.Mul(t, new(big.Int).Set(s))
	t.ModInverse(t, q)
	t.Mul(r, t)
	t.Mod(t, q)

	u := big.NewInt(1 << bias)
	u.Mul(u, s)
	u.Mod(u, q)
	u.ModInverse(u, q)
	u.Mul(z, u)
	u.Mod(u, q)
	u.Sub(q, u)

	return t, u
}

func main() {
	f, err := os.Create("/tmp/profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	p, _ := new(big.Int).SetString("233970423115425145524320034830162017933", 10)
	a := big.NewInt(-95051)
	b := big.NewInt(11279326)
	G := dh.NewEllipticCurve(a, b, p)

	gx := big.NewInt(182)
	gy, _ := new(big.Int).SetString("85518893674295321206118380980485522083", 10)
	g := dh.NewEllipticCurveElement(G, gx, gy)
	q, _ := new(big.Int).SetString("29246302889428143187362802287225875743", 10)
	GG := dh.NewGeneratedGroup(G, g, *q)

	var bias uint = 8
	alice := dh.NewBiasedECDSA(GG, bias)

	var msg = []byte("I was a fiend")
	key, _ := new(big.Int).SetString("255bf9c75628ab469b45cced58755a3", 16)
	d := new(big.Rat).SetInt(key)

	numSigs := 22
	B := dh.Matrix(make([]dh.Vector, numSigs+2))
	zero := dh.Vector(make([]*big.Rat, numSigs+2))
	for i, _ := range zero {
		zero[i] = new(big.Rat)
	}
	for i, _ := range B {
		B[i] = zero.Copy()
	}
	ct := big.NewRat(1, 1<<bias)
	B[len(B)-2][len(B)-2].Set(ct)
	cu := new(big.Rat).SetInt(q)
	cu.Quo(cu, big.NewRat(1<<bias, 1))
	B[len(B)-1][len(B)-1].Set(cu)
	ts := make([]*big.Int, numSigs)
	us := make([]*big.Int, numSigs)
	for i := 0; i < numSigs; i++ {
		B[i][i].SetInt(q)
		r, s := alice.Sign(msg)
		t, u := transform(msg, r, s, q, bias)
		dt := new(big.Int).Mul(key, t)
		temp := new(big.Int).Sub(u, dt)
		temp.Mod(temp, q)
		temp.Sub(q, temp)
		log.Printf("\ndt:     %x\nu:      %x\nq-u-dt: %x\nq:      %x", dt, u, temp, q)
		ts[i] = t
		us[i] = u
		B[len(B)-2][i] = new(big.Rat).SetInt(t)
		B[len(B)-1][i] = new(big.Rat).SetInt(u)
	}

	check := B[len(B)-1].Copy()
	check.Sub(B[len(B)-2].Copy().Scale(d))
	for i := 0; i < numSigs; i++ {
		t := ts[i]
		dt := new(big.Int).Mul(key, t)
		m := new(big.Int).Div(dt, q)
		check.Add(B[i].Copy().Scale(new(big.Rat).SetInt(m)))
	}
	log.Printf("check=%s", check)

	B.LLL(big.NewRat(99, 100))

	for _, v := range B {
		if v[len(v)-1].Cmp(cu) == 0 {
			log.Printf("%s", v)
			d := new(big.Rat)
			d.Sub(d, v[len(v)-2])
			d.Mul(d, big.NewRat(1<<bias, 1))
			guess := dh.Round(d).Num()
			log.Printf("Recovered key: %x", guess)
			log.Printf("Correct: %v", guess.Cmp(key) == 0)
		}
	}
}
