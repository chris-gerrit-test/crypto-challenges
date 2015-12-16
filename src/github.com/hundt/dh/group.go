package dh

import (
	"fmt"
	"math/big"
	"math/rand"
)

var rnd = rand.New(rand.NewSource(77778))

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
	Size() big.Int
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

func (gg *generatedGroup) Size() big.Int {
	return *gg.q
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

func (e *ellipticCurveElement) String() string {
	if e == inf {
		return "∞"
	}
	return fmt.Sprintf("(%s,%s)", e.x, e.y)
}

func (e *ellipticCurveElement) Jump(k int) int {
	panic("not implemented")
	return -1
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
			denom.ModInverse(denom, ec.modulus)
			m.Mul(m, denom)
			m.Mod(m, ec.modulus)
		}
	}
	if m == nil {
		m = new(big.Int)
		denom := new(big.Int)
		denom.Sub(e2.x, e1.x)
		denom.Mod(denom, ec.modulus)
		denom.ModInverse(denom, ec.modulus)
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
