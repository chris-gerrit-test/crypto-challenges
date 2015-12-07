package dh

import (
	"fmt"
	"math/big"
	"math/rand"
)

type Group interface {
	Op(x, y GroupMember) GroupMember
	Pow(x GroupMember, y *big.Int) GroupMember
	Random() GroupMember
}

type CyclicGroup interface {
	Group
	Generator() GroupMember
	Size() big.Int
}

type GroupMember interface {
	Group() Group
	String() string
}

type finiteGroup struct {
	p *big.Int
}

func NewFiniteGroup(p big.Int) Group {
	return &finiteGroup{p: &p}
}

type finiteGroupMember struct {
	n *big.Int
	g *finiteGroup
}

func NewFiniteGroupMember(g Group, n big.Int) GroupMember {
	pg := g.(*finiteGroup)
	return &finiteGroupMember{g: pg, n: n.Mod(&n, pg.p)}
}

func (n *finiteGroupMember) Group() Group {
	return n.g
}

func (n *finiteGroupMember) String() string {
	return n.n.String()
}

func (pg *finiteGroup) Op(x, y GroupMember) GroupMember {
	px := x.(*finiteGroupMember)
	py := y.(*finiteGroupMember)
	z := new(big.Int)
	z.Mul(px.n, py.n)
	z.Mod(z, pg.p)
	return &finiteGroupMember{n: z, g: pg}
}

func (pg *finiteGroup) Pow(x GroupMember, y *big.Int) GroupMember {
	px := x.(*finiteGroupMember)
	z := new(big.Int)
	z.Exp(px.n, y, pg.p)
	return &finiteGroupMember{n: z, g: pg}
}

var rnd = rand.New(rand.NewSource(77777))

func (pg *finiteGroup) Random() GroupMember {
	z := new(big.Int)
	for z.Cmp(new(big.Int)) == 0 {
		z.Rand(rnd, pg.p)
	}
	return &finiteGroupMember{n: z, g: pg}
}

type generatedGroup struct {
	g Group
	m GroupMember
	q *big.Int // order of g
}

type generatedGroupMember struct {
	m  GroupMember
	gg *generatedGroup
}

func (m *generatedGroupMember) Group() Group {
	return m.gg
}

func (m *generatedGroupMember) String() string {
	return m.m.String()
}

func NewGeneratedGroup(m GroupMember, q big.Int) CyclicGroup { // q is order of g in mod n
	return &generatedGroup{g: m.Group(), m: m, q: &q}
}

func (gg *generatedGroup) String() string {
	return fmt.Sprintf("{size %s generated from %s}", gg.q, gg.Generator())
}

func (gg *generatedGroup) Op(x, y GroupMember) GroupMember {
	return &generatedGroupMember{
		m:  gg.g.Op(x.(*generatedGroupMember).m, y.(*generatedGroupMember).m),
		gg: gg,
	}
}

func (gg *generatedGroup) Pow(x GroupMember, y *big.Int) GroupMember {
	return &generatedGroupMember{
		m:  gg.g.Pow(x.(*generatedGroupMember).m, y),
		gg: gg,
	}
}

func (gg *generatedGroup) Random() GroupMember {
	z := new(big.Int)
	z.Rand(rnd, gg.q)
	return &generatedGroupMember{
		m:  gg.g.Pow(gg.m, z),
		gg: gg,
	}
}

func (gg *generatedGroup) Generator() GroupMember {
	return &generatedGroupMember{
		m:  gg.m,
		gg: gg,
	}
}

func (gg *generatedGroup) Size() big.Int {
	return *gg.q
}
