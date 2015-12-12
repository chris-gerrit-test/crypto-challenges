package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"log"
	"math/big"
	"math/rand"

	"github.com/hundt/dh"
)

func Generate(ch chan<- int64) {
	for i := int64(2); ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int64, out chan<- int64, prime int64) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// Sends primes to the passed-in channel forever
func StreamPrimes(ch chan<- int64) {
	ch2 := make(chan int64) // Create a new channel.
	go Generate(ch2)        // Launch Generate goroutine.
	for {
		prime := <-ch2
		ch <- prime
		ch1 := make(chan int64)
		go Filter(ch2, ch1, prime)
		ch2 = ch1
	}
}

var secretMessage = []byte("crazy flamboyant for the rap enjoyment")

func Sign(msg []byte, key dh.Element) []byte {
	h := hmac.New(sha1.New, []byte(key.String()))
	return h.Sum(msg)
}

func Bob(d dh.DiffieHellman, pub dh.Element) []byte {
	k := d.SharedSecret(pub)
	return Sign(secretMessage, k)
}

var rnd = rand.New(rand.NewSource(9825))

func main() {
	pi := new(big.Int)
	pi.SetString("7199773997391911030609999317773941274322764333428698921736339643928346453700085358802973900485592910475480089726140708102474957429903531369589969318716771", 10)
	gi := new(big.Int)
	gi.SetString("4565356397095740655436854503483826832136106141639563487732438195343690437606117828318042418238184896212352329118608100083187535033402010599512641674644143", 10)
	q := new(big.Int)
	q.SetString("236234353446506858198510045061214171961", 10)

	j := new(big.Int)
	j.Sub(pi, big.NewInt(1))
	j.Div(j, q)

	ch := make(chan int64) // Create a new channel.
	go StreamPrimes(ch)
	zero := new(big.Int)
	factors := make(map[int64]*big.Int, 0)
	total := big.NewInt(1)
	for total.Cmp(q) < 0 {
		prime := <-ch
		pr := big.NewInt(prime)
		pr.Rem(j, pr)
		if pr.Cmp(zero) == 0 {
			j2 := new(big.Int).Set(j)
			j2.Div(j, big.NewInt(prime))
			pr.Rem(j2, big.NewInt(prime))
			if pr.Cmp(zero) == 0 {
				log.Printf("Skipping double divisor %d", prime)
			} else {
				log.Printf("Found divisor %d", prime)
				factors[prime] = nil
				total.Mul(total, big.NewInt(prime))
			}
		}
	}

	for factor, _ := range factors {
		h := new(big.Int)
		groupSize := new(big.Int)
		groupSize.Sub(pi, big.NewInt(1))
		pow := new(big.Int)
		pow.Div(groupSize, big.NewInt(factor))
		log.Printf("Finding an element of order %d...", factor)
		for {
			h.Rand(rnd, pi)
			h.Exp(h, pow, pi)
			if h.Cmp(big.NewInt(1)) != 0 {
				break
			}
		}
		factors[factor] = h
	}

	G := dh.NewFiniteGroup(*pi)
	g := dh.NewFiniteElement(G, *gi)
	GG := dh.NewGeneratedGroup(G, g, *q)
	d := dh.NewDiffieHellman(GG)

	moduli := make(map[int64]int64)

	for factor, elt := range factors {
		h := dh.NewFiniteElement(G, *elt)
		mac := Bob(d, h)
		e := new(big.Int).Set(elt)
		log.Printf("Guessing the shared secret in the subgroup of order %d", factor)
		found := false
		for i := int64(1); i <= factor; i++ {
			e.Exp(elt, big.NewInt(i), pi)
			k := dh.NewFiniteElement(G, *e)
			if hmac.Equal(mac, Sign(secretMessage, k)) {
				//log.Printf("%d^%d", elt, i)
				found = true
				moduli[factor] = i
				break
			}
		}
		if !found {
			panic("Could not guess the shared secret")
		}
	}

	// From Wikipedia CRT page
	x := new(big.Int)
	for n, a := range moduli {
		N := new(big.Int).Set(total)
		N.Div(N, big.NewInt(n))
		N.ModInverse(N, big.NewInt(n))
		N.Mul(N, total)
		N.Div(N, big.NewInt(n))
		N.Mul(N, big.NewInt(a))
		x.Add(x, N)
		x.Mod(x, total)
	}
	log.Printf("Predicted key: %d", x)

	log.Printf("%s", d)
}
