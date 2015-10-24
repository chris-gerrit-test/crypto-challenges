package main

import (
	"hash"
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

type Node struct {
	state *TruncatedMD4
	next  *Node
	msg   []byte // msg takes state -> next.state
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

func PaddingBlock(byteLength int) []byte {
	b := make([]byte, blockSize)
	b[0] = 0x80

	bitLength := byteLength << 3
	for i := uint(0); i < 8; i++ {
		b[blockSize-8+i] = byte(bitLength >> (8 * i))
	}

	return b
}

func main() {
	k := 12
	leaves := make([]*Node, 0)

	msg := getMessage(1)

	for i := 0; i < 1<<uint(k); i++ {
		lastState := NewTruncatedMD4()
		lastState.MustWrite(msg)
		leaves = append(leaves, &Node{state: lastState})
		nextMessage(msg)
	}

	current := leaves
	depth := 0
	for len(current) > 1 {
		next := make([]*Node, len(current)/2)
		log.Printf("Finding %d collisions at depth %d", len(next), depth)
		for i := 0; i < len(next); i++ {
			n1 := current[2*i]
			n2 := current[2*i+1]
			//log.Printf("%x %x", n1.state.SumNoPad(nil), n2.state.SumNoPad(nil))
			n1Hashes := make(map[string][]byte, 0) // nextState -> msg
			n2Hashes := make(map[string][]byte, 0)
			msg := getMessage(1)
			var s *TruncatedMD4
			for {
				s = n1.state.Copy()
				s.MustWrite(msg)
				sum := string(s.SumNoPad(nil))
				if b, ok := n2Hashes[sum]; ok {
					//log.Printf("Found match in n2Hashes: %x\n", sum)
					n1.msg = msg
					n2.msg = b
					break
				}
				n1Hashes[sum] = append([]byte(nil), msg...)
				s = n2.state.Copy()
				s.MustWrite(msg)
				sum = string(s.SumNoPad(nil))
				if b, ok := n1Hashes[sum]; ok {
					//log.Printf("Found match in n1Hashes: %x\n", sum)
					n1.msg = b
					n2.msg = msg
					break
				}
				n2Hashes[sum] = append([]byte(nil), msg...)
				nextMessage(msg)
			}
			next[i] = &Node{state: s}
			n1.next = next[i]
			n2.next = next[i]
		}
		depth++
		current = next
	}

	// for _, n := range leaves {
	// 	s := n.state.Copy()
	// 	msg := make([]byte, 0)
	// 	for n.next != nil {
	// 		msg = append(msg, n.msg...)
	// 		n = n.next
	// 	}
	// 	s.MustWrite(msg)
	// 	log.Printf("%x: %s", s.SumNoPad(nil), msg)
	// 	log.Printf("%x (padded sum)", s.Sum(nil))
	// 	s.MustWrite(PaddingBlock(len(msg) + blockSize)) // one block to get initial collision
	// 	log.Printf("%x (my padded sum)", s.SumNoPad(nil))
	// }

	predictionBlocks := 4
	finalState := current[0].state.Copy()
	finalMessageLength := blockSize * (predictionBlocks + k + 1)
	log.Printf("The final message will be %d bytes long with %d bytes of prediction",
		finalMessageLength, blockSize*predictionBlocks)
	finalState.MustWrite(PaddingBlock(finalMessageLength))
	finalMessageHashPredicted := finalState.SumNoPad(nil)
	log.Printf("Final message hash will be %x", finalMessageHashPredicted)

	prediction := `The Toronto Blue Jays will make it to the ALCS in 2015, where they will lose to the Kansas City Royals.

In the ninth inning there will be two questionable called strikes against the Blue Jays.`

	predictionPadding := strings.Repeat(" ", predictionBlocks*blockSize-len(prediction))
	log.Printf("Prediction is %d bytes long. Adding %d bytes of padding", len(prediction), len(predictionPadding))
	prediction += predictionPadding
	if len(prediction)%blockSize != 0 || len(prediction)/blockSize != predictionBlocks {
		panic(len(prediction))
	}

	leafHashes := make(map[string]int, 0)
	for i, n := range leaves {
		leafHashes[string(n.state.SumNoPad(nil))] = i
	}
	log.Printf("Finding a collision into a leaf")
	pState := NewTruncatedMD4()
	pState.MustWrite([]byte(prediction))
	glue := getMessage(1)
	finalMessage := ""
	for {
		s := pState.Copy()
		s.MustWrite(glue)
		if i, ok := leafHashes[string(s.SumNoPad(nil))]; ok {
			//log.Printf("Found match in leafHashes: %d/%x\n", i, s.SumNoPad(nil))
			finalMessage = prediction + string(glue)
			n := leaves[i]
			for n.next != nil {
				finalMessage += string(n.msg)
				n = n.next
			}
			break
		}
		nextMessage(glue)
	}

	if len(finalMessage) != finalMessageLength {
		panic(len(finalMessage))
	}
	pState = NewTruncatedMD4()
	pState.MustWrite([]byte(finalMessage))
	log.Printf("Final message:\n%s", finalMessage)
	log.Printf("Hash: %x", pState.Sum(nil))
}
