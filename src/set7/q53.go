package main

import (
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hundt/md4"
)

const hashSize = 5
const blockSize = 64

type TruncatedMD4 struct {
	hash.Hash
	len int
}

func NewTruncatedMD4() *TruncatedMD4 {
	return &TruncatedMD4{md4.New(hashSize), hashSize}
}

func (t *TruncatedMD4) Sum(b []byte) []byte {
	return t.Hash.Sum(b)[:t.len]
}

func (t *TruncatedMD4) SumNoPad(b []byte) []byte {
	return md4.SumNoPad(t.Hash, b)[:t.len]
}

func (t *TruncatedMD4) Copy() *TruncatedMD4 {
	return &TruncatedMD4{md4.Copy(t.Hash), t.len}
}

func (t *TruncatedMD4) MustWrite(b []byte) {
	n, err := t.Write(b)
	if n != len(b) || err != nil {
		panic(err)
	}
}

type Collision struct {
	small, big []byte
	state      *TruncatedMD4
}

func getMessage(blocks int) []byte {
	return []byte(strings.Repeat("!", blocks*blockSize))
}

func nextMessage(msg []byte) []byte {
	for i := len(msg) - 1; i >= 0; i-- {
		msg[i] += 1
		if msg[i] <= '~' {
			break
		}
		msg[i] = '!'
	}
	return msg
}

func main() {
	k := 20
	collisions := make([]*Collision, k)
	lastState := NewTruncatedMD4()
	for i := 1; i <= k; i++ {
		log.Printf("Looking for a collision with %d blocks\n", 1+(1<<uint(k-i)))
		c := &Collision{state: lastState.Copy(), small: getMessage(1), big: getMessage(1 + (1 << uint(k-i)))}
		bigState := c.state.Copy()
		bigState.MustWrite(c.big[:len(c.big)-blockSize])
		smalls := make(map[string][]byte, 0)
		bigs := make(map[string][]byte, 0)
		for {
			nextMessage(c.small)
			newState := c.state.Copy()
			newState.MustWrite(c.small)
			sum := string(newState.SumNoPad(nil))
			if b, ok := bigs[sum]; ok {
				//fmt.Printf("Found match in bigs: %x\n", sum)
				copy(c.big[len(c.big)-blockSize:], b)
				c.state = newState
				break
			}
			smalls[sum] = append([]byte(nil), c.small...)

			nextMessage(c.big[len(c.big)-blockSize:])
			newState = bigState.Copy()
			newState.MustWrite(c.big[len(c.big)-blockSize:])
			sum = string(newState.SumNoPad(nil))
			if b, ok := smalls[sum]; ok {
				//fmt.Printf("Found match in smalls: %x\n", sum)
				c.small = b
				c.state = newState
				break
			}
			bigs[sum] = append([]byte(nil), c.big[len(c.big)-blockSize:]...)
		}
		bc := lastState.Copy()
		bc.MustWrite(c.big)
		sc := lastState.Copy()
		sc.MustWrite(c.small)
		// fmt.Printf("State: %x\n", c.state.SumNoPad(nil))
		// fmt.Printf("Big:\n  %s\n  (%x)\nSmall:\n  %s\n  (%x)\n",
		// 	c.big, bc.SumNoPad(nil), c.small, sc.SumNoPad(nil))
		lastState = c.state
		collisions[i-1] = c
	}

	h := NewTruncatedMD4()
	for i, collision := range collisions {
		if i%2 == 0 {
			h.MustWrite(collision.big)
		} else {
			h.MustWrite(collision.small)
		}
	}
	// log.Printf("Mixed:      %x", h.SumNoPad(nil))
	// log.Printf("Last state: %x", collisions[len(collisions)-1].state.SumNoPad(nil))

	m := getMessage(1 << uint(k))
	h = NewTruncatedMD4()
	h.MustWrite(m)
	log.Printf("Original message: %d bytes", len(m))
	log.Printf("Original message hash: %x\n", h.Sum(nil))
	ioutil.WriteFile("/tmp/original.txt", m, 0644)

	log.Printf("Computing intermediate hashes")
	intermediates := make(map[string]int, 0) // intermediate hash -> block index
	intermediate := NewTruncatedMD4()
	for bi := 0; bi < 1<<uint(k); bi++ {
		if bi >= k+1 {
			intermediates[string(intermediate.SumNoPad(nil))] = bi
		}
		intermediate.MustWrite(m[bi*blockSize : (1+bi)*blockSize])
	}

	log.Printf("Finding a collision into the intermediates")
	bridge := getMessage(1)
	startState := collisions[len(collisions)-1].state
	matchIndex := -1
	for {
		s := startState.Copy()
		s.MustWrite(bridge)
		sum := string(s.SumNoPad(nil))
		if bi, ok := intermediates[sum]; ok {
			matchIndex = bi
			//log.Printf("Found collision: %x", sum)
			h := NewTruncatedMD4()
			h.MustWrite(m[:bi*blockSize])
			//log.Printf("Message prefix sum: %x", h.SumNoPad(nil))
			break
		}
		nextMessage(bridge)
	}

	prefixBlocks := matchIndex - 1
	//log.Printf("Need %d prefix blocks", prefixBlocks)
	forged := make([]byte, 0)
	for i, collision := range collisions {
		var piece []byte
		if prefixBlocks >= len(collision.big)/blockSize+len(collisions)-i-1 {
			piece = collision.big
		} else {
			piece = collision.small
		}
		forged = append(forged, piece...)
		prefixBlocks -= len(piece) / blockSize
	}
	if prefixBlocks != 0 {
		panic(fmt.Sprintf("There are still %d blocks remaining", prefixBlocks))
	}
	// h = NewTruncatedMD4()
	// h.MustWrite(forged)
	// log.Printf("Prefix hash: %x", h.SumNoPad(nil))
	// h.MustWrite(bridge)
	// log.Printf("With bridge: %x", h.SumNoPad(nil))

	forged = append(forged, bridge...)
	forged = append(forged, m[matchIndex*blockSize:]...)
	log.Printf("Forged message: %d bytes", len(forged))
	h = NewTruncatedMD4()
	h.MustWrite(forged)
	log.Printf("Forged message hash: %x", h.Sum(nil))

	ioutil.WriteFile("/tmp/forged.txt", forged, 0644)
}
