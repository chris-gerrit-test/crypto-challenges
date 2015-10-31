package main

import (
	"encoding/hex"
	"log"
	"strings"
)

func ToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var shift1 = []uint{3, 7, 11, 19}
var shift2 = []uint{3, 5, 9, 13}
var shift3 = []uint{3, 9, 11, 15}

var xIndex2 = []uint{0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14, 3, 7, 11, 15}
var xIndex3 = []uint{0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15}

func CheckConditions(m []byte) {
	a, b, c, d := uint32(0x67452301), uint32(0xEFCDAB89), uint32(0x98BADCFE), uint32(0x10325476)

	var X [16]uint32

	j := 0
	for i := 0; i < 16; i++ {
		X[i] = uint32(m[j]) | uint32(m[j+1])<<8 | uint32(m[j+2])<<16 | uint32(m[j+3])<<24
		j += 4
	}

	// Round 1.
	b0 := b
	for i := uint(0); i < 16; i++ {
		x := i
		s := shift1[i%4]
		f := ((c ^ d) & b) ^ d
		a += f + X[x]
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c
		if i == 3 {
			log.Printf("Condition a1: %t", a&0x40 == b0&0x40)
			cond := (d&0x40 == 0) && (d&0x80 == a&0x80) && (d&0x400 == a&0x400)
			log.Printf("Condition d1: %t", cond)
			cond = (c&0x40 == 0x40) && (c&0x80 == 0x80) && (c&0x400 == 0) && (c&0x2000000 == d&0x2000000)
			log.Printf("Condition c1: %t", cond)
			cond = (b&0x40 == 0x40) && (b&0x80 == 0) && (b&0x400 == 0) && (b&0x2000000 == 0)
			log.Printf("Condition b1: %t", cond)
		}
	}

	c4 := c

	// Round 2.
	for i := uint(0); i < 16; i++ {
		x := xIndex2[i]
		s := shift2[i%4]
		g := (b & c) | (b & d) | (c & d)
		a += g + X[x] + 0x5a827999
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c

		if i == 3 {
			cond := ((a&0x40000 == c4&0x40000) &&
				(a&0x2000000 == 0x2000000) &&
				(a&0x4000000 == 0) &&
				(a&0x10000000 == 0x10000000) &&
				(a&0x80000000 == 0x80000000))
			log.Printf("Condition a5: %t", cond)
		}
	}

	// Round 3.
	for i := uint(0); i < 16; i++ {
		x := xIndex3[i]
		s := shift3[i%4]
		h := b ^ c ^ d
		a += h + X[x] + 0x6ed9eba1
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c
	}

}

func Correct(m []byte) {
	a, b, c, d := uint32(0x67452301), uint32(0xEFCDAB89), uint32(0x98BADCFE), uint32(0x10325476)
	a0 := a

	var X [16]uint32
	var Xp [16]uint32
	var Xc [16]uint32

	j := 0
	for i := 0; i < 16; i++ {
		X[i] = uint32(m[j]) | uint32(m[j+1])<<8 | uint32(m[j+2])<<16 | uint32(m[j+3])<<24
		Xp[i] = X[i]
		Xc[i] = X[i]
		j += 4
	}
	Xp[1] = X[1] + 0x80000000
	Xp[2] = X[2] + 0x80000000 - 0x10000000
	Xp[12] = X[12] - 0x10000

	// Round 1.
	for i := uint(0); i < 16; i++ {
		x := i
		s := shift1[i%4]
		f := ((c ^ d) & b) ^ d
		a += f + X[x]
		a = a<<s | a>>(32-s)
		if i == 0 {
			// Correct for a1 constraint
			ac := a ^ ((a & 0x40) ^ (b & 0x40))
			Xc[0] = (ac>>s | ac<<(32-s)) - a0 - f
		}
		a, b, c, d = d, a, b, c
	}

	//c4 := c

    func correct(i uint32) {
        
    }

	// Round 2.
	for i := uint(0); i < 16; i++ {
		x := xIndex2[i]
		s := shift2[i%4]
		g := (b & c) | (b & d) | (c & d)
		a += g + X[x] + 0x5a827999
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c

		if i == 3 {
		}
	}

	// Round 3.
	for i := uint(0); i < 16; i++ {
		x := xIndex3[i]
		s := shift3[i%4]
		h := b ^ c ^ d
		a += h + X[x] + 0x6ed9eba1
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c
	}

	m[0] = byte(Xc[0])
	m[1] = byte(Xc[0] >> 8)
	m[2] = byte(Xc[0] >> 16)
	m[3] = byte(Xc[0] >> 24)
}

func main() {
	// m1 := ToBytes("4d7a9c8356cb927ab9d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dd8e3197e31fe52794bf08b9e8c3e9")
	// //m1p := ToBytes("4d7a9c83d6cb927a29d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dc8e3197e31fe52794bf08b9e8c3e9")
	// m2 := ToBytes("4d7a9c8356cb927ab9d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dd8e3197e31fe5f713c240a7b8cf69")
	// for i := 0; i < len(m1); i += 4 {
	// 	for j := 0; j < 2; j++ {
	// 		i1 := i + j
	// 		i2 := i + 3 - j
	// 		c := m1[i1]
	// 		m1[i1] = m1[i2]
	// 		m1[i2] = c
	// 		c = m2[i1]
	// 		m2[i1] = m2[i2]
	// 		m2[i2] = c
	// 	}
	// }
	// CheckConditions(m1)
	// CheckConditions(m2)
	// md1 := md4.New(16)
	// md2 := md4.New(16)
	// // md1.Write(m1)
	// // md2.Write(m2)
	// // log.Printf("%x", md1.Sum(nil))
	// // log.Printf("%x", md2.Sum(nil))

	m := []byte(strings.Repeat("z", 64))
	CheckConditions(m)

	Correct(m)
	CheckConditions(m)
}
