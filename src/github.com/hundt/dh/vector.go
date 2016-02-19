package dh

import (
	"log"
	"math/big"
)

type Vector []*big.Rat

type Matrix []Vector

func (m Matrix) String() string {
	s := "("
	for _, v := range m {
		s += "\n "
		s += v.String()
	}
	return s + "\n)"
}

func (m1 Matrix) Copy() Matrix {
	m2 := Matrix(make([]Vector, len(m1)))
	for i, v := range m1 {
		m2[i] = v.Copy()
	}
	return m2
}

func (m1 Matrix) Eq(m2 Matrix) bool {
	for i, v := range m1 {
		if !v.Eq(m2[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) GramSchmidt() Matrix {
	if len(m) == 0 {
		return m
	}
	p := m[0].Copy()
	for i, vi := range m {
		for j := i + 1; j < len(m); j++ {
			vj := m[j]
			p.Set(vj).Proj(vi)
			vj.Sub(p)
		}
	}
	return m
}

func mu(b, q Matrix, i, j int) *big.Rat {
	v := b[i]
	u := q[j]
	d := v.Dot(u)
	return d.Quo(d, u.Dot(u))
}

func Round(r *big.Rat) *big.Rat {
	d := new(big.Int).Set(r.Denom())
	n := new(big.Int).Set(r.Num())
	n.Mod(n, d)
	if new(big.Int).Mul(n, big.NewInt(2)).Cmp(d) >= 0 {
		r.Add(r, new(big.Rat).SetInt64(1))
	}
	r.Sub(r, new(big.Rat).SetFrac(n, d))
	return r
}

func (b Matrix) LLL(delta *big.Rat) Matrix {
	q := b.Copy().GramSchmidt()

	n := len(b)
	k := 1

	for k < n {
		for j := k - 1; j >= 0; j-- {
			m := mu(b, q, k, j)
			m.Abs(m)
			if m.Cmp(big.NewRat(1, 2)) > 0 {
				v := b[j].Copy()
				m := mu(b, q, k, j)
				b[k].Sub(v.Scale(Round(m)))
				q = b.Copy().GramSchmidt()
			}
		}

		lhs := q[k].Dot(q[k])
		rhs := q[k-1].Dot(q[k-1])
		m := mu(b, q, k, k-1)
		m.Mul(m, m)
		m.Sub(delta, m)
		rhs.Mul(m, rhs)
		if lhs.Cmp(rhs) >= 0 {
			k += 1
		} else {
			for i := 0; i < len(b[k]); i++ {
				temp := new(big.Rat).Set(b[k][i])
				b[k][i].Set(b[k-1][i])
				b[k-1][i].Set(temp)
			}
			q = b.Copy().GramSchmidt()
			k -= 1
			if k < 1 {
				k = 1
			}
		}
	}

	return b
}

func (v1 Vector) Add(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		e.Add(e, v2[i])
	}
	return v1
}

func (v1 Vector) Sub(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		e.Sub(e, v2[i])
	}
	return v1
}

func (v1 Vector) Eq(v2 Vector) bool {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		if e.Cmp(v2[i]) != 0 {
			return false
		}
	}
	return true
}

func (v1 Vector) Dot(v2 Vector) *big.Rat {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	d := new(big.Rat)
	c := new(big.Rat)
	for i, e := range v1 {
		d.Add(d, c.Mul(e, v2[i]))
	}
	return d
}

func (v1 Vector) Copy() Vector {
	v2 := Vector(make([]*big.Rat, len(v1)))
	for i, e := range v1 {
		v2[i] = new(big.Rat).Set(e)
	}
	return v2
}

func (v1 Vector) Set(v2 Vector) Vector {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	for i, e := range v1 {
		e.Set(v2[i])
	}
	return v1
}

func (v Vector) String() string {
	s := "("
	for i, e := range v {
		if i != 0 {
			s += ", "
		}
		s += e.String()
	}
	return s + ")"
}

func (v Vector) Scale(r *big.Rat) Vector {
	for _, e := range v {
		e.Mul(e, r)
	}
	return v
}

func (v1 Vector) Proj(v2 Vector) Vector {
	isZero := true
	for _, e := range v2 {
		if e.Cmp(big.NewRat(0, 1)) != 0 {
			isZero = false
			break
		}
	}
	if isZero {
		v1.Set(v2)
	} else {
		u1 := v1.Copy()
		u2 := v2.Copy()
		d1 := u1.Dot(v2)
		d1.Quo(d1, u2.Dot(v2))
		v1.Set(v2).Scale(d1)
	}
	return v1
}
