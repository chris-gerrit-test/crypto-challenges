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

func (gg *generatedGroup) Generator() Element {
	return gg.m
}

func (gg *generatedGroup) Size() big.Int {
	return *gg.q
}
