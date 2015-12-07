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
	PublicKey() GroupMember
	SharedSecret(B GroupMember) GroupMember
}

func (d *diffieHellman) String() string {
	return fmt.Sprintf("{DH group=%s key=%s}", d.g, d.a)
}

func NewDiffieHellman(g CyclicGroup) DiffieHellman {
	z := new(big.Int)
	n := g.Size()
	z.Rand(rnd, &n)
	return &diffieHellman{g: g, a: z}
}

func (d *diffieHellman) PublicKey() GroupMember {
	return d.g.Pow(d.g.Generator(), d.a)
}

func (d *diffieHellman) SharedSecret(B GroupMember) GroupMember {
	return d.g.Pow(B, d.a)
}
