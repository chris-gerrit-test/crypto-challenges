package main

import (
	"log"
	"math/big"

	"github.com/hundt/dh"
)

func main() {
	g := dh.NewFiniteGroup(*big.NewInt(211))
	x := dh.NewFiniteGroupMember(g, *big.NewInt(5))
	gg := dh.NewGeneratedGroup(g, x, *big.NewInt(35))
	//y := dh.NewPrimeFiniteGroupMember(p, *big.NewInt(2))
	seen := make(map[string]int)
	for i := 0; i < 100000; i++ {
		seen[gg.Random().String()] += 1
	}
	for n, count := range seen {
		log.Printf("%3d: %d", n, count)
	}
	log.Printf("%d entries", len(seen))
}
