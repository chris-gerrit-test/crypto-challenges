package dh

import (
	"log"
	"math/big"
)

type Vector []*big.Rat

type Matrix []*Vector

func (v1 Vector) Add(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		e.Add(e, v2[i])
	}
	return v1
}

func (v1 Vector) Eq(v2 Vector) bool {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		if e.Cmp(v2[i]) != 0 {
			return false
		}
	}
	return true
}

func (v1 Vector) Dot(v2 Vector) *big.Rat {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	d := new(big.Rat)
	c := new(big.Rat)
	for i, e := range v1 {
		d.Add(d, c.Mul(e, v2[i]))
	}
	return d
}

func (v1 Vector) Copy() Vector {
	v2 := Vector(make([]*big.Rat, len(v1)))
	for i, e := range v1 {
		v2[i] = new(big.Rat).Set(e)
	}
	return v2
}

func (v1 Vector) Set(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		e.Set(v2[i])
	}
	return v1
}

func (v Vector) String() string {
	s := "("
	for i, e := range v {
		if i != 0 {
			s += ", "
		}
		s += e.String()
	}
	return s + ")"
}

func (v Vector) Scale(r *big.Rat) Vector {
	for _, e := range v {
		e.Mul(e, r)
	}
	return v
}

func (v1 Vector) Proj(v2 Vector) Vector {
	u1 := v1.Copy()
	u2 := v2.Copy()
	d1 := u1.Dot(v2)
	d1.Quo(d1, u2.Dot(v2))
	v1.Set(v2).Scale(d1)
	return v1
}
