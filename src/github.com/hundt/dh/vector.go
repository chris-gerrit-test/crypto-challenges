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

func (m Matrix) fixAfterUpdated(updated int) Matrix {
	if len(m) == 0 {
		return m
	}
	p := m[0].Copy() // temp storage
	vj := m[updated]
	for i := 0; i < updated; i++ {
		vi := m[i]
		p.Set(vj).Proj(vi)
		vj.Sub(p)
	}
	// Adjust vectors after the updated ones
	for i := updated; i < len(m)-1; i++ {
		vi := m[i]
		norm := vi.fastDot(vi)
		for j := i + 1; j < len(m); j++ {
			vj := m[j]
			p.Set(vj).fastProj(vi, norm)
			vj.Sub(p)
		}
	}
	return m
}

func (m Matrix) fixUpdated(first, last int) Matrix {
	if len(m) == 0 {
		return m
	}
	p := m[0].Copy() // temp storage
	for i := 0; i < first; i++ {
		vi := m[i]
		norm := vi.fastDot(vi)

		vj := m[first]
		p.Set(vj).fastProj(vi, norm)
		vj.Sub(p)

		vj = m[last]
		p.Set(vj).fastProj(vi, norm)
		vj.Sub(p)
	}
	vi := m[first]
	vj := m[last]
	p.Set(vj).Proj(vi)
	vj.Sub(p)
	return m
}

func mu(b, q Matrix, i, j int) *big.Rat {
	v := b[i]
	u := q[j]
	d := v.fastDot(u)
	return d.Quo(d, u.fastDot(u))
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
				// q2 := b.Copy().GramSchmidt()
				// if !q.Eq(q2) {
				// 	log.Fatalf("Already not orthogonal")
				// }
				v := b[j].Copy()
				m := mu(b, q, k, j)
				b[k].Sub(v.Scale(Round(m)))
				q[k].Set(b[k].Copy())
				q.fixAfterUpdated(k)
				// q2 = b.Copy().GramSchmidt()
				// if !q.Eq(q2) {
				// 	log.Fatalf("you fucked up A: (k=%d)\n%s\n%s", k, q, q2)
				// }
			}
		}

		lhs := q[k].fastDot(q[k])
		rhs := q[k-1].fastDot(q[k-1])
		m := mu(b, q, k, k-1)
		m.Mul(m, m)
		m.Sub(delta, m)
		rhs.Mul(m, rhs)
		if lhs.Cmp(rhs) >= 0 {
			k += 1
		} else {
			// q2 := b.Copy().GramSchmidt()
			// if !q.Eq(q2) {
			// 	log.Fatalf("Already not orthogonal")
			// }
			temp := b[k].Copy()
			b[k].Set(b[k-1])
			b[k-1].Set(temp)
			q[k].Set(b[k])
			q[k-1].Set(b[k-1])
			q.fixUpdated(k-1, k)
			// q2 = b.Copy().GramSchmidt()
			// if !q.Eq(q2) {
			// 	log.Fatalf("you fucked up B: (k=%d)\n%s\n%s", k, q, q2)
			// }
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

func (v Vector) totalDenom() *big.Int {
	d := big.NewInt(1)
	for _, e := range v {
		d.Mul(d, e.Denom())
	}
	return d
}

func (v1 Vector) fastDot(v2 Vector) *big.Rat {
	if len(v1) != len(v2) {
		log.Fatalf("Lengths differ: %d != %d", len(v1), len(v2))
	}
	denom := v1.totalDenom()
	denom.Mul(denom, v2.totalDenom())
	temp := new(big.Int)
	num := new(big.Int)
	for i, e := range v1 {
		temp.Mul(e.Denom(), v2[i].Denom())
		temp.Div(denom, temp)
		temp.Mul(temp, e.Num())
		temp.Mul(temp, v2[i].Num())
		num.Add(num, temp)
	}
	return new(big.Rat).SetFrac(num, denom)
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

func (v1 Vector) fastProj(v2 Vector, v2Norm *big.Rat) Vector {
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
		d1 := u1.fastDot(v2)
		d1.Quo(d1, v2Norm)
		v1.Set(v2).Scale(d1)
	}
	return v1
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
		d1 := u1.fastDot(v2)
		d1.Quo(d1, u2.fastDot(v2))
		v1.Set(v2).Scale(d1)
	}
	return v1
}
