package dh

import (
	"fmt"
	"log"
	"math/big"
)

type polyElement struct {
	f  Field
	cs []Element
}

type polyField struct {
	p *polyElement
}

func NewPolyField(p Element) Field {
	return &polyField{p: p.(*polyElement)}
}

func NewPolyElement(f Field, cs []Element) Element {
	return &polyElement{f: f, cs: cs}
}

func (f *polyField) Identity() Element {
	return &finiteElement{n: big.NewInt(1)}
}

func (f *polyField) Zero() Element {
	log.Fatal("not implemented")
	return nil
}

func (f *polyField) Op(x, y Element) Element {
	log.Fatal("not implemented")
	return nil
}

func mod(e *polyElement, f *polyField) {
	// for i, c := range e.cs {
	// 	mc := f.p.cs[i]
	// 	if c.Cmp(mc) >= 0 {
	// 		e.cs[i] = e.f.Sub(c, mc)
	// 	}
	// }
}

func (f *polyField) Add(x, y Element) Element {
	p1 := x.(*polyElement)
	p2 := y.(*polyElement)
	sums := make([]Element, len(p1.cs))
	for i, c1 := range p1.cs {
		c2 := p2.cs[i]
		sums[i] = p1.f.Add(c1, c2)
	}
	sum := &polyElement{f: p1.f, cs: sums}
	mod(sum, f)
	return sum
}

func (f *polyField) Sub(x, y Element) Element {
	p1 := x.(*polyElement)
	p2 := y.(*polyElement)
	sums := make([]Element, len(p1.cs))
	for i, c1 := range p1.cs {
		c2 := p2.cs[i]
		sums[i] = p1.f.Sub(c1, c2)
	}
	sum := &polyElement{f: p1.f, cs: sums}
	mod(sum, f)
	return sum
}

func (f *polyField) Pow(x Element, y *big.Int) Element {
	log.Fatal("not implemented")
	return nil
}

func (f *polyField) Random() Element {
	log.Fatal("not implemented")
	return nil
}

func (p *polyElement) String() string {
	s := ""
	for i := len(p.cs) - 1; i >= 0; i-- {
		c := p.cs[i]
		if c.Cmp(p.f.Zero()) != 0 {
			if s != "" {
				s += " + "
			}
			if i == 1 {
				s += fmt.Sprintf("%s*x", c)
			} else if i == 0 {
				s += fmt.Sprintf("%s", c)
			} else {
				s += fmt.Sprintf("%s*x^%d", c, i)
			}
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

func (p *polyElement) Jump(k int) int {
	log.Panic("not implemented")
	return 0
}

func (p *polyElement) Cmp(e Element) int {
	p2 := e.(*polyElement)
	for i := len(p.cs) - 1; i >= 0; i-- {
		c1 := p.cs[i]
		c2 := p2.cs[i]
		cmp := c1.Cmp(c2)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}
