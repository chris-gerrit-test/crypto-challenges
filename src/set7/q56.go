package main

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/base64"
	"log"
	"strings"
)

type Peephole struct {
	request          []byte
	counts1, counts2 map[byte]int
}

func main() {
	counts := make([]map[byte]int, 32)
	for i, _ := range counts {
		counts[i] = make(map[byte]int, 0)
	}
	secret, err := base64.StdEncoding.DecodeString("QkUgU1VSRSBUTyBEUklOSyBZT1VSIE9WQUxUSU5F")
	if err != nil || len(secret) != 30 {
		panic(err)
	}
	peepholes := make(map[int]*Peephole)
	for off := -1; off < 15; off++ {
		peepholes[off] = &Peephole{
			request: []byte("/" + strings.Repeat("x", 14-off) + string(secret)),
			counts1: make(map[byte]int, 0),
			counts2: make(map[byte]int, 0),
		}
	}
	keySize := 8
	key := make([]byte, keySize)

	s := make([]byte, len(secret)+16)
	for trial := 0; trial < 10000000; trial++ {
		for idx, p := range peepholes {
			n, err := rand.Read(key)
			if n != keySize || err != nil {
				panic(err)
			}
			c, err := rc4.NewCipher(key)
			if err != nil {
				panic(err)
			}
			c.XORKeyStream(s, p.request)
			p.counts1[s[15]^240] += 1
			if idx < 15 {
				p.counts2[s[31]^224] += 1
			}
		}

		if trial%100000 == 99999 {
			// Show our best guess so far
			for idx, p := range peepholes {
				if idx >= 0 {
					s[idx] = '?'
					max := 0
					for b, count := range p.counts1 {
						// Assume ASCII without any funny whitespace
						if b >= 32 && b <= 126 && count > max {
							s[idx] = b
							max = count
						}
					}
				}
				if idx < 15 {
					s[idx+16] = '?'
					max := 0
					for b, count := range p.counts2 {
						if b >= 32 && b <= 126 && count > max {
							s[idx+16] = b
							max = count
						}
					}
				}
			}
			log.Printf("Best guess after %8d trials: %s", trial+1, s[:30])
		}
	}
}
