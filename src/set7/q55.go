package main

import (
	"bytes"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hundt/md4"
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

func CheckConditions(m []byte) bool {
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
			//log.Printf("Condition a1: %t", a&0x40 == b0&0x40)
			if a&0x40 != b0&0x40 {
				log.Printf("Condition a1: %t", a&0x40 == b0&0x40)
				return false
			}
			cond := (d&0x40 == 0) && (d&0x80 == a&0x80) && (d&0x400 == a&0x400)
			if !cond {
				log.Printf("Condition d1: %t", cond)
				return false
			}
			//log.Printf("Condition d1: %t", cond)
			cond = (c&0x40 == 0x40) && (c&0x80 == 0x80) && (c&0x400 == 0) && (c&0x2000000 == d&0x2000000)
			if !cond {
				log.Printf("Condition c1: %t", cond)
				return false
			}
			//log.Printf("Condition c1: %t", cond)
			// cond = (b&0x40 == 0x40) && (b&0x80 == 0) && (b&0x400 == 0) && (b&0x2000000 == 0)
			//log.Printf("Condition b1: %t", cond)
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
			// log.Printf("Condition a5.1: %t", a&0x40000 == c4&0x40000)
			// log.Printf("Condition a5.2: %t", a&0x2000000 == 0x2000000)
			// log.Printf("Condition a5.3: %t", a&0x4000000 == 0)
			// log.Printf("Condition a5.4: %t", a&0x10000000 == 0x10000000)
			// log.Printf("Condition a5.5: %t", a&0x80000000 == 0x80000000)
			// log.Printf("Condition a5: %t", cond)
			if !cond {
				return false
			}
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

	return true

}

func rrot(n uint32, r uint) uint32 {
	return (n >> r) | n<<(32-r)
}

func Correct(m []byte, cond uint32) {
	a, b, c, d := uint32(0x67452301), uint32(0xEFCDAB89), uint32(0x98BADCFE), uint32(0x10325476)

	var X [16]uint32
	var Xc [16]uint32

	j := 0
	for i := 0; i < 16; i++ {
		X[i] = uint32(m[j]) | uint32(m[j+1])<<8 | uint32(m[j+2])<<16 | uint32(m[j+3])<<24
		Xc[i] = X[i]
		j += 4
	}

	var A [9]uint32
	var B [9]uint32
	var C [9]uint32
	var D [9]uint32

	A[0] = a
	B[0] = b
	C[0] = c
	D[0] = d

	F := func(b, c, d uint32) uint32 {
		return ((c ^ d) & b) ^ d
	}

	// Round 1.
	for i := uint(0); i < 16; i++ {
		x := i
		s := shift1[i%4]
		f := ((c ^ d) & b) ^ d
		a += f + X[x]
		a = a<<s | a>>(32-s)

		a, b, c, d = d, a, b, c
		if i%4 == 3 {
			A[i/4+1] = a
			B[i/4+1] = b
			C[i/4+1] = c
			D[i/4+1] = d
		}
	}
	if cond&0x1 != 0 {
		//log.Printf("Correct a1")
		ac := A[1] ^ ((A[1] & 0x40) ^ (B[0] & 0x40))
		Xc[0] = rrot(ac, shift1[0]) - A[0] - F(B[0], C[0], D[0])
	}
	if cond&0x2 != 0 {
		dc := D[1] ^ (D[1] & 0x40) ^ ((D[1] ^ A[1]) & 0x80) ^ ((D[1] ^ A[1]) & 0x400)
		Xc[1] = rrot(dc, 7) - D[0] - F(A[1], B[0], C[0])
	}

	// Round 2.
	for i := uint(0); i < 16; i++ {
		x := xIndex2[i]
		s := shift2[i%4]
		g := (b & c) | (b & d) | (c & d)
		a += g + X[x] + 0x5a827999
		a = a<<s | a>>(32-s)
		a, b, c, d = d, a, b, c
		if i%4 == 3 {
			A[i/4+5] = a
			B[i/4+5] = b
			C[i/4+5] = c
			D[i/4+5] = d
		}
	}
	// Correct for a5 constraints
	type Correction struct {
		i   uint32
		add bool
	}
	corrections := make([]Correction, 0)
	aNew := A[1]
	if cond&0x4 != 0 && A[5]&0x40000 != C[4]&0x40000 {
		//log.Printf("Correct a5.1")
		corrections = append(corrections, Correction{19, C[4]&0x40000 != 0})
	}
	if cond&0x8 != 0 && A[5]&0x2000000 != 0x2000000 {
		//log.Printf("Correct a5.2")
		corrections = append(corrections, Correction{26, true})
	}
	if cond&0x10 != 0 && A[5]&0x4000000 != 0 {
		//log.Printf("Correct a5.3")
		corrections = append(corrections, Correction{27, false})
	}
	if cond&0x20 != 0 && A[5]&0x10000000 != 0x10000000 {
		//log.Printf("Correct a5.4")
		corrections = append(corrections, Correction{29, true})
	}
	if cond&0x40 != 0 && A[5]&0x80000000 != 0x80000000 {
		//log.Printf("Correct a5.5")
		corrections = append(corrections, Correction{32, true})
	}
	for _, corr := range corrections {
		adj := uint32(1) << (corr.i - 4)
		mask := uint32(1) << (corr.i - 1)
		if aNew&mask == 0 {
			Xc[0] += adj
			aNew = aNew | mask
		} else {
			Xc[0] -= adj
			aNew = aNew &^ mask
		}
		Xc[1] = rrot(D[1], 7) - D[0] - F(aNew, B[0], C[0])
		Xc[2] = rrot(C[1], 11) - C[0] - F(D[1], aNew, B[0])
		Xc[3] = rrot(B[1], 19) - B[0] - F(C[1], D[1], aNew)
		Xc[4] = rrot(A[2], 3) - aNew - F(B[1], C[1], D[1])
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

	for i := 0; i < 5; i++ {
		m[4*i] = byte(Xc[i])
		m[4*i+1] = byte(Xc[i] >> 8)
		m[4*i+2] = byte(Xc[i] >> 16)
		m[4*i+3] = byte(Xc[i] >> 24)
	}
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

func GetPair(M, Mp []byte) {
	var X [16]uint32
	var Xp [16]uint32

	j := 0
	for i := 0; i < 16; i++ {
		X[i] = uint32(M[j]) | uint32(M[j+1])<<8 | uint32(M[j+2])<<16 | uint32(M[j+3])<<24
		Xp[i] = X[i]
		j += 4
	}
	Xp[1] = X[1] + 0x80000000
	Xp[2] = X[2] + 0x80000000 - 0x10000000
	Xp[12] = X[12] - 0x10000
	for i := 0; i < 16; i++ {
		Mp[4*i] = byte(Xp[i])
		Mp[4*i+1] = byte(Xp[i] >> 8)
		Mp[4*i+2] = byte(Xp[i] >> 16)
		Mp[4*i+3] = byte(Xp[i] >> 24)
	}
}

func main() {
	m1 := ToBytes("4d7a9c8356cb927ab9d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dd8e3197e31fe52794bf08b9e8c3e9")
	m1p := ToBytes("4d7a9c83d6cb927a29d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dc8e3197e31fe52794bf08b9e8c3e9")
	m2 := ToBytes("4d7a9c8356cb927ab9d5a57857a7a5eede748a3cdcc366b3b683a0203b2a5d9fc69d71b3f9e99198d79f805ea63bb2e845dd8e3197e31fe5f713c240a7b8cf69")
	// Swap little/big endian
	for i := 0; i < len(m1); i += 4 {
		for j := 0; j < 2; j++ {
			i1 := i + j
			i2 := i + 3 - j
			c := m1[i1]
			m1[i1] = m1[i2]
			m1[i2] = c
			c = m2[i1]
			m2[i1] = m2[i2]
			m2[i2] = c
			c = m1p[i1]
			m1p[i1] = m1p[i2]
			m1p[i2] = c
		}
	}
	m1pp := bytes.Repeat([]byte{'x'}, len(m1))
	GetPair(m1, m1pp)
	// Verify that my Pair method works
	if string(m1pp) != string(m1p) {
		log.Printf("%x", m1pp)
		log.Printf("%x", m1p)
		return
	}
	// CheckConditions(m1)
	// CheckConditions(m2)
	// md1 := md4.New(16)
	// md2 := md4.New(16)
	// // md1.Write(m1)
	// // md2.Write(m2)
	// // log.Printf("%x", md1.Sum(nil))
	// // log.Printf("%x", md2.Sum(nil))

	// m := []byte(strings.Repeat("x", 64))
	// log.Printf("%x", m)
	// CheckConditions(m)
	// Correct(m, 1)
	// Correct(m, 2)
	// Correct(m, 4)
	// Correct(m, 8)
	// Correct(m, 16)
	// Correct(m, 32)
	// log.Printf("%x", m)
	// CheckConditions(m)

	M0 := []byte(strings.Repeat("a", 64))
	M := make([]byte, len(M0))
	Mp := []byte(strings.Repeat("x", 64))
	for trials := 1; ; trials++ {
		copy(M, M0)
		for i := uint32(1); i <= 64; i *= 2 {
			Correct(M, i)
		}
		if !CheckConditions(M) {
			return
		}
		GetPair(M, Mp)
		dig := md4.New(16)
		dig.Write(M)
		m := dig.Sum(nil)
		dig = md4.New(16)
		dig.Write(Mp)
		mp := dig.Sum(nil)
		if string(m) == string(mp) {
			log.Printf("%x: %x", m, M)
			log.Printf("%x: %x", mp, Mp)
			log.Printf("%d tries", trials)
			return
		}
		if trials%1000000 == 0 {
			log.Printf("%d tries", trials)
		}
		nextMessage(M0)
	}

}
