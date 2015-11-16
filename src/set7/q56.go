package main

import (
	"bufio"
	"crypto/rc4"
	"fmt"
	"os"
)

type Count struct {
	b     byte
	count uint64
}

type Counts []Count

func (c Counts) Len() int {
	return len(c)
}

func (c Counts) Less(i, j int) bool {
	return c[i].count < c[j].count
}

func (c Counts) Swap(i, j int) {
	z := c[i]
	c[i] = c[j]
	c[j] = z
}

func main() {
	counts := make([]map[byte]uint64, 32)
	for i, _ := range counts {
		counts[i] = make(map[byte]uint64, 0)
	}
	secret := "YELLOW SUBMARINEYELLOW SUBMARINE"
	s := make([]byte, len(secret))
	keySize := 8
	key := make([]byte, keySize)

	f, err := os.Open("random.dat")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReaderSize(f, 16384)
	var off int64 = 0
	for i := 0; i < 100000000; i++ {
		copy(s, secret)
		//n, err := rand.Read(key)
		n, err := r.Read(key)
		if n != keySize || err != nil {
			panic(err)
		}
		off += int64(keySize)
		c, err := rc4.NewCipher(key)
		if err != nil {
			panic(err)
		}
		c.XORKeyStream(s, s)
		for b := 0; b < len(counts); b++ {
			counts[b][s[b]^secret[b]] += 1
		}
	}

	// for c, count := range counts {
	// 	log.Printf("%3d: %d", c^byte(240), count)
	// }

	for b := 0; b < len(counts); b++ {
		for c := 0; c < 256; c++ {
			fmt.Printf("C%02d %3d %d\n", b, c, counts[b][byte(c)])
		}
	}
}
