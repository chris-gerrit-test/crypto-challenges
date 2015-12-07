package main

import (
	"log"
	"math/big"

	"github.com/hundt/dh"
)

func main() {
	g := dh.NewFiniteGroup(*big.NewInt(211))
	x := dh.NewFiniteGroupMember(g, *big.NewInt(5))
	gg := dh.NewGeneratedGroup(x, *big.NewInt(35))
	dh1 := dh.NewDiffieHellman(gg)
	dh2 := dh.NewDiffieHellman(gg)
	log.Printf("\n%s\n%s", dh1, dh2)
	log.Printf("%s %s", dh1.PublicKey(), dh2.PublicKey())
	log.Printf("%s %s", dh1.SharedSecret(dh2.PublicKey()), dh2.SharedSecret(dh1.PublicKey()))
}
