package dh

import (
	"fmt"
	"math/big"
)

type diffieHellman struct {
	g CyclicGroup
	a *big.Int // a == g^n
}

type DiffieHellman interface {
	PublicKey() Element
	SharedSecret(B Element) Element
}

func (d *diffieHellman) String() string {
	return fmt.Sprintf("{DH group=%s key=%s}", d.g, d.a)
}

func NewDiffieHellman(g CyclicGroup) DiffieHellman {
	z := new(big.Int)
	z.Rand(rnd, g.Size())
	return &diffieHellman{g: g, a: z}
}

func (d *diffieHellman) PublicKey() Element {
	return d.g.Pow(d.g.Generator(), d.a)
}

func (d *diffieHellman) SharedSecret(B Element) Element {
	return d.g.Pow(B, d.a)
}
