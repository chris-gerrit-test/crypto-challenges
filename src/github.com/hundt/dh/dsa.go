package dh

import (
	"crypto/sha1"
	"log"
	"math/big"
)

type ecdsa struct {
	g   CyclicGroup
	key *big.Int
}

type ECDSA interface {
	Group() CyclicGroup
	PublicKey() Element
	Sign(m []byte) (*big.Int, *big.Int)
}

// the generator of g should come from NewEllipticCurveElement
func NewECDSA(g CyclicGroup) ECDSA {
	n := g.Size()
	d := &ecdsa{g: g, key: new(big.Int)}
	for d.key.Cmp(new(big.Int)) == 0 {
		d.key.Rand(rnd, n)
	}
	return d
}

func (d *ecdsa) PublicKey() Element {
	return d.g.Pow(d.g.Generator(), d.key)
}

func (d *ecdsa) Group() CyclicGroup {
	return d.g
}

func (d *ecdsa) Sign(m []byte) (*big.Int, *big.Int) {
	h := sha1.New()
	if n, err := h.Write(m); n != len(m) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	r, s := new(big.Int), new(big.Int)
	n := d.g.Size()
	z := new(big.Int).SetBytes(e)
	z.Mod(z, n)
	for r.Cmp(new(big.Int)) == 0 || s.Cmp(new(big.Int)) == 0 {
		k := new(big.Int).Rand(rnd, n)
		if k.Cmp(new(big.Int)) == 0 {
			continue
		}
		p := d.g.Pow(d.g.Generator(), k)
		r.Mod(p.(*ellipticCurveElement).x, n)

		k.ModInverse(k, n)
		s.Mul(r, d.key)
		s.Add(s, z)
		s.Mul(s, k)
		s.Mod(s, n)
	}
	return r, s
}

func ECDSAVerify(m []byte, r, s *big.Int, d ECDSA) bool {
	Q := d.PublicKey()
	g := d.Group()
	n := g.Size()
	if s.Cmp(new(big.Int)) <= 0 || s.Cmp(n) >= 0 {
		return false
	}
	if r.Cmp(new(big.Int)) <= 0 || r.Cmp(n) >= 0 {
		return false
	}
	h := sha1.New()
	if n, err := h.Write(m); n != len(m) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	z := new(big.Int).SetBytes(e)
	z.Mod(z, n)
	w := new(big.Int).ModInverse(s, n)
	u1 := new(big.Int).Mul(z, w)
	u1.Mod(u1, n)
	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, n)
	p1 := g.Pow(g.Generator(), u1)
	p2 := g.Pow(Q, u2)
	R := g.Op(p1, p2).(*ellipticCurveElement)
	x := R.x
	x.Mod(x, n)
	return x.Cmp(r) == 0
}

// Returns an ECDSA instance for which Sign(m) == r, s and which will
// be in a subgroup of the same size as g.
func FindVerifyingECDSA(m []byte, r, s *big.Int, g CyclicGroup) ECDSA {
	n := g.Size()
	d := &ecdsa{g: g, key: new(big.Int)}
	for d.key.Cmp(new(big.Int)) == 0 {
		d.key.Rand(rnd, n)
	}
	h := sha1.New()
	if n, err := h.Write(m); n != len(m) || err != nil {
		log.Fatal("Error calculating hash")
	}
	e := h.Sum(nil)
	z := new(big.Int).SetBytes(e)
	z.Mod(z, n)
	w := new(big.Int).ModInverse(s, n)
	u1 := new(big.Int).Mul(z, w)
	u1.Mod(u1, n)
	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, n)
	rr := new(big.Int).Set(r)
	t := new(big.Int).Mul(u2, d.key)
	t.Add(t, u1)
	t.ModInverse(t, n)
	// If we don't have the public key that actually generated the signature,
	// we need to guess what R was before doing r = R.x mod n
	for ; ; rr.Add(rr, g.Size()) {
		R, err := NewEllipticCurveElementFromX(g.(*generatedGroup).g, rr)
		if err != nil {
			continue
		}
		for _, invert := range []bool{false, true} {
			if invert {
				NegateWeierstrass(R, g.(*generatedGroup).g)
			}
			G := g.Pow(R, t)
			d.g = NewGeneratedGroup(g.(*generatedGroup).g, G, *d.g.Size())
			if ECDSAVerify(m, r, s, d) {
				return d
			}
		}
	}
}
