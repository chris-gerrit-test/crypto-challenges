package dh

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
)

var rnd = rand.New(rand.NewSource(77780))

func jump(x *big.Int, k int) int {
	j := new(big.Int)
	j.Mod(x, big.NewInt(int64(k)))
	return 1 << j.Uint64()
}

type Group interface {
	Op(x, y Element) Element
	Pow(x Element, y *big.Int) Element
	Random() Element
	Identity() Element
}

type CyclicGroup interface {
	Group
	Generator() Element
	Size() *big.Int
}

type Element interface {
	String() string
	Cmp(Element) int
	Jump(k int) int
}

type finiteGroup struct {
	p *big.Int
}

func NewFiniteGroup(p big.Int) Group {
	return &finiteGroup{p: &p}
}

type finiteElement struct {
	n *big.Int
}

func NewFiniteElement(g Group, n big.Int) Element {
	pg := g.(*finiteGroup)
	return &finiteElement{n: n.Mod(&n, pg.p)}
}

func (n *finiteElement) String() string {
	return n.n.String()
}

func (n *finiteElement) Jump(k int) int {
	return jump(n.n, k)
}

func (n *finiteElement) Cmp(e Element) int {
	m := e.(*finiteElement)
	return n.n.Cmp(m.n)
}

func (pg *finiteGroup) Identity() Element {
	return &finiteElement{n: big.NewInt(1)}
}

func (pg *finiteGroup) Op(x, y Element) Element {
	px := x.(*finiteElement)
	py := y.(*finiteElement)
	z := new(big.Int)
	z.Mul(px.n, py.n)
	z.Mod(z, pg.p)
	return &finiteElement{n: z}
}

func (pg *finiteGroup) Pow(x Element, y *big.Int) Element {
	px := x.(*finiteElement)
	z := new(big.Int)
	z.Exp(px.n, y, pg.p)
	return &finiteElement{n: z}
}

func (pg *finiteGroup) Random() Element {
	z := new(big.Int)
	for z.Cmp(new(big.Int)) == 0 {
		z.Rand(rnd, pg.p)
	}
	return &finiteElement{n: z}
}

type generatedGroup struct {
	g Group
	m Element
	q *big.Int // order of m in g
}

func NewGeneratedGroup(g Group, m Element, q big.Int) CyclicGroup { // q is order of m in g
	return &generatedGroup{m: m, q: &q, g: g}
}

func (gg *generatedGroup) String() string {
	return fmt.Sprintf("{size %s generated from %s}", gg.q, gg.Generator())
}

func (gg *generatedGroup) Op(x, y Element) Element {
	return gg.g.Op(x, y)
}

func (gg *generatedGroup) Pow(x Element, y *big.Int) Element {
	return gg.g.Pow(x, y)
}

func (gg *generatedGroup) Random() Element {
	z := new(big.Int)
	z.Rand(rnd, gg.q)
	return gg.g.Pow(gg.m, z)
}

func (gg *generatedGroup) Identity() Element {
	return gg.g.Identity()
}

func (gg *generatedGroup) Generator() Element {
	return gg.m
}

func (gg *generatedGroup) Size() *big.Int {
	return new(big.Int).Set(gg.q)
}

type ellipticCurve struct {
	a, b    *big.Int
	modulus *big.Int
}

func (ec *ellipticCurve) String() string {
	return fmt.Sprintf("y^2 = x^3 + %s*x + %s (mod %s)", ec.a, ec.b, ec.modulus)
}

func NewEllipticCurve(a, b, modulus *big.Int) Group {
	return &ellipticCurve{
		a:       new(big.Int).Set(a),
		b:       new(big.Int).Set(b),
		modulus: new(big.Int).Set(modulus),
	}
}

type ellipticCurveElement struct {
	x, y *big.Int
}

var inf *ellipticCurveElement = new(ellipticCurveElement)

func NewEllipticCurveElement(g Group, x, y *big.Int) Element {
	// TODO: check that the points are on the curve?
	// TODO mod
	return &ellipticCurveElement{
		x: new(big.Int).Set(x),
		y: new(big.Int).Set(y),
	}
}

func NewEllipticCurveElementFromX(g Group, x *big.Int) (Element, error) {
	ec := g.(*ellipticCurve)
	x3 := new(big.Int).Exp(x, big.NewInt(3), ec.modulus)
	y := new(big.Int).Mul(ec.a, x)
	y.Add(y, x3)
	y.Add(y, ec.b)
	if y.ModSqrt(y, ec.modulus) == nil {
		return nil, errors.New("No square root")
	}
	return &ellipticCurveElement{
		x: new(big.Int).Set(x),
		y: new(big.Int).Set(y),
	}, nil
}

func (e *ellipticCurveElement) String() string {
	if e == inf {
		return "∞"
	}
	return fmt.Sprintf("(%s,%s)", e.x, e.y)
}

func (e *ellipticCurveElement) Jump(k int) int {
	return jump(e.x, k)
}

func (e1 *ellipticCurveElement) Cmp(e Element) int {
	e2 := e.(*ellipticCurveElement)
	if e1 == inf {
		if e2 == inf {
			return 0
		}
		return 1
	}
	if e2 == inf {
		return -1
	}
	c := e1.x.Cmp(e2.x)
	if c != 0 {
		return c
	}
	return e1.y.Cmp(e2.y)
}

func (ec *ellipticCurve) Op(x, y Element) Element {
	e1 := x.(*ellipticCurveElement)
	e2 := y.(*ellipticCurveElement)
	if e1 == inf {
		return e2
	}
	if e2 == inf {
		return e1
	}
	pMinus2 := new(big.Int).Sub(ec.modulus, big.NewInt(2))
	var m *big.Int
	if e1.x.Cmp(e2.x) == 0 {
		my2 := new(big.Int).Set(e2.y)
		my2.Sub(ec.modulus, my2)
		if my2.Cmp(e1.y) == 0 {
			return inf
		}
		if e2.y.Cmp(e1.y) == 0 {
			m = new(big.Int)
			m.Exp(e1.x, big.NewInt(2), ec.modulus)
			m.Mul(m, big.NewInt(3))
			m.Add(m, ec.a)
			denom := new(big.Int)
			denom.Mul(big.NewInt(2), e1.y)
			denom.Exp(denom, pMinus2, ec.modulus)
			//denom.ModInverse(denom, ec.modulus)
			m.Mul(m, denom)
			m.Mod(m, ec.modulus)
		}
	}
	if m == nil {
		m = new(big.Int)
		denom := new(big.Int)
		denom.Sub(e2.x, e1.x)
		denom.Mod(denom, ec.modulus)
		denom.Exp(denom, pMinus2, ec.modulus)
		//denom.ModInverse(denom, ec.modulus)
		m.Sub(e2.y, e1.y)
		m.Mul(m, denom)
		m.Mod(m, ec.modulus)
	}
	e := &ellipticCurveElement{
		x: new(big.Int),
		y: new(big.Int),
	}
	e.x.Mul(m, m)
	e.x.Sub(e.x, e1.x)
	e.x.Sub(e.x, e2.x)
	e.x.Mod(e.x, ec.modulus)
	e.y.Sub(e1.x, e.x)
	e.y.Mul(m, e.y)
	e.y.Sub(e.y, e1.y)
	e.y.Mod(e.y, ec.modulus)
	return e
}

func genericPow(G Group, x Element, y *big.Int) Element {
	y = new(big.Int).Set(y)
	var result Element = nil
	zero := new(big.Int)
	bit := new(big.Int)
	two := big.NewInt(2)
	if y.Cmp(zero) < 0 {
		panic("Cannot use genericPow with y < 0")
	}
	for y.Cmp(zero) != 0 {
		y.QuoRem(y, two, bit)
		if bit.Cmp(zero) != 0 {
			if result == nil {
				result = x
			} else {
				result = G.Op(result, x)
			}
		}
		x = G.Op(x, x)
	}
	return result
}

func (pg *ellipticCurve) Pow(x Element, y *big.Int) Element {
	return genericPow(pg, x, y)
}

func (pg *ellipticCurve) Random() Element {
	x := new(big.Int)
	var y *big.Int = nil
	t := new(big.Int)
	for y == nil {
		x.Rand(rnd, pg.modulus)
		xc := new(big.Int).Exp(x, big.NewInt(3), pg.modulus)
		t.Mul(x, pg.a)
		t.Add(t, xc)
		t.Add(t, pg.b)
		t.Mod(t, pg.modulus)
		y = t.ModSqrt(t, pg.modulus)
	}
	return &ellipticCurveElement{
		x: x,
		y: y,
	}
}

func (gg *ellipticCurve) Identity() Element {
	return inf
}

type montgomeryCurve struct {
	A, B    *big.Int
	modulus *big.Int
}

func (ec *montgomeryCurve) String() string {
	return fmt.Sprintf("%s*v^2 = u^3 + %s*u^2 + u (mod %s)", ec.B, ec.A, ec.modulus)
}

func NegateWeierstrass(e Element, g Group) {
	elt := e.(*ellipticCurveElement)
	grp := g.(*ellipticCurve)
	elt.y.Sub(grp.modulus, elt.y)
}

func Q60MontgomeryToWeierstrass(e Element, newGroup Group) Element {
	g := newGroup.(*ellipticCurve)
	elt := e.(*montgomeryCurveElement)
	x := new(big.Int).Add(elt.u, big.NewInt(178))
	sc1 := new(big.Int).Exp(elt.u, big.NewInt(3), g.modulus)
	sc2 := new(big.Int).Mul(elt.u, elt.u)
	sc2.Mul(sc2, big.NewInt(534))
	sc1.Add(sc1, sc2)
	sc1.Add(sc1, elt.u)
	if sc1.ModSqrt(sc1, g.modulus) == nil {
		log.Panic("no matching y found")
	}
	return NewEllipticCurveElement(g, x, sc1)
}

func Q60WeierstrassToMontgomery(x Element, newGroup Group) Element {
	e := x.(*ellipticCurveElement)
	u := new(big.Int).Set(e.x)
	u.Sub(u, big.NewInt(178))
	u.Mod(u, newGroup.(*montgomeryCurve).modulus)
	return NewMontgomeryCurveElement(newGroup, u)
}

func NewMontgomeryCurve(A, B, modulus *big.Int) Group {
	return &montgomeryCurve{
		A:       new(big.Int).Set(A),
		B:       new(big.Int).Set(B),
		modulus: new(big.Int).Set(modulus),
	}
}

type montgomeryCurveElement struct {
	u *big.Int
}

func NewMontgomeryCurveElement(g Group, u *big.Int) Element {
	return &montgomeryCurveElement{
		u: new(big.Int).Set(u),
	}
}

func (e *montgomeryCurveElement) String() string {
	if e.u.Cmp(new(big.Int)) == 0 {
		return "∞"
	}
	return fmt.Sprintf("%s", e.u)
}

func (e *montgomeryCurveElement) Jump(k int) int {
	panic("not implemented")
	return -1
}

func (e1 *montgomeryCurveElement) Cmp(e Element) int {
	e2 := e.(*montgomeryCurveElement)
	return e1.u.Cmp(e2.u)
}

func (ec *montgomeryCurve) Op(x, y Element) Element {
	log.Fatal("Not implemented")
	return nil
}

func cswap(x, y *big.Int, b uint) (*big.Int, *big.Int) {
	if b == 1 {
		return y, x
	}
	return x, y
}

func (pg *montgomeryCurve) Pow(x Element, k *big.Int) Element {
	e := x.(*montgomeryCurveElement)
	u2 := big.NewInt(1)
	w2 := big.NewInt(0)
	u3 := new(big.Int).Set(e.u)
	w3 := big.NewInt(1)
	sc1 := new(big.Int)
	sc2 := new(big.Int)
	sc3 := new(big.Int)
	foundBit := false
	for i := pg.modulus.BitLen() - 1; i >= 0; i-- {
		b := k.Bit(i)
		if foundBit || b != 0 {
			foundBit = true
		}
		if !foundBit {
			continue
		}
		//log.Printf("Begin: (%d, %d)    (%d, %d)", u2, w2, u3, w3)
		//log.Printf("bit %d is %d", i, b)
		u2, u3 = cswap(u2, u3, b)
		w2, w3 = cswap(w2, w3, b)
		//log.Printf("After cswap: (%d, %d)    (%d, %d)", u2, w2, u3, w3)
		// u3 = u2*u3 - w2*w3)^2
		sc1.Mul(u2, u3)
		sc2.Mul(w2, w3)
		sc1.Sub(sc1, sc2)
		sc3.Mul(sc1, sc1)
		sc3.Mod(sc3, pg.modulus)
		// w3 = u1 * (u2*w3 - w2*u3)^2
		sc1.Mul(u2, w3)
		sc2.Mul(w2, u3)
		sc1.Sub(sc1, sc2)
		w3.Mul(sc1, sc1)
		w3.Mul(e.u, w3)
		w3.Mod(w3, pg.modulus)
		u3.Set(sc3)
		// u2 = (u2^2 - w2^2)^2
		sc1.Mul(u2, u2)
		sc2.Mul(w2, w2)
		sc3.Sub(sc1, sc2)
		sc3.Mul(sc3, sc3)
		sc3.Mod(sc3, pg.modulus)
		// w2 = 4*u2*w2 * (u2^2 + A*u2*w2 + w2^2)
		sc1.Mul(u2, u2)
		sc2.Mul(pg.A, u2)
		sc2.Mul(sc2, w2)
		sc1.Add(sc1, sc2)
		sc2.Mul(w2, w2)
		sc1.Add(sc1, sc2)
		sc2.Mul(big.NewInt(4), u2)
		sc2.Mul(sc2, w2)
		w2.Mul(sc2, sc1)
		w2.Mod(w2, pg.modulus)
		u2.Set(sc3)

		u2, u3 = cswap(u2, u3, b)
		w2, w3 = cswap(w2, w3, b)
		//log.Printf("End: (%d, %d)    (%d, %d)", u2, w2, u3, w3)
	}
	//log.Printf("%d %d", u2, w2)
	sc1.Sub(pg.modulus, big.NewInt(2))
	sc1.Exp(w2, sc1, pg.modulus)
	sc1.Mul(sc1, u2)
	sc1.Mod(sc1, pg.modulus)
	return &montgomeryCurveElement{
		u: sc1,
	}
}

// NB: this actually returns a random point that is *not* on the curve.
func (pg *montgomeryCurve) Random() Element {
	if pg.B.Cmp(big.NewInt(1)) != 0 {
		log.Fatal("Only B=1 is supported")
	}
	for {
		u := new(big.Int).Rand(rnd, pg.modulus)
		square := new(big.Int).Mul(u, u)
		square.Mul(square, pg.A)
		cube := new(big.Int).Exp(u, big.NewInt(3), pg.modulus)
		v := new(big.Int).Add(square, cube)
		v.Add(v, u)
		v.Mod(v, pg.modulus)
		if v.ModSqrt(v, pg.modulus) == nil {
			return &montgomeryCurveElement{
				u: u,
			}
		}
	}
}

func (gg *montgomeryCurve) Identity() Element {
	return &montgomeryCurveElement{
		u: new(big.Int),
	}
}
