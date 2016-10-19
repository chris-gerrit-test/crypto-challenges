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
	log.Panic("not implemented")
	return nil
}

func (f *polyField) Inverse(e Element) Element {
	log.Fatal("not implemented")
	return nil
}

func (f *polyField) Zero() Element {
	sums := make([]Element, len(f.p.cs)-1)
	for i, _ := range sums {
		sums[i] = f.p.f.Zero()
	}
	return NewPolyElement(f.p.f, sums).(*polyElement)
}

func (f *polyField) Op(x, y Element) Element {
	p1 := x.(*polyElement)
	p2 := y.(*polyElement)
	// make copies
	p1 = f.Sub(p1, f.Zero()).(*polyElement)
	p2 = f.Sub(p2, f.Zero()).(*polyElement)
	p := f.Zero().(*polyElement)

	zero := p1.f.Zero()
	for _ = range p1.cs {
		log.Printf("%d %d %d", len(p2.cs), len(p1.cs), len(f.p.cs))
		log.Printf("\na = %s\nb = %s", p1, p2)
		for i, c2 := range p2.cs {
			p.cs[i] = p.f.Add(p.cs[i], p.f.Op(p1.cs[0], c2))
		}

		p1.cs = append(p1.cs[1:len(p1.cs)], p1.f.Zero())
		p2.cs = append([]Element{p2.f.Zero()}, p2.cs...)
		log.Printf("\na = %s\nb = %s", p1, p2)

		hi := p2.cs[len(p2.cs)-1]
		if hi.Cmp(zero) != 0 {
			log.Printf("modding %s", p2)
			for i, c2 := range p2.cs {
				p2.cs[i] = p2.f.Sub(c2, p2.f.Op(hi, f.p.cs[i]))
			}
			log.Printf("got %s", p2)
		}
		if p2.cs[len(p2.cs)-1].Cmp(zero) != 0 {
			panic("modding didn't work")
		}

		p2.cs = p2.cs[0 : len(p2.cs)-1]
	}

	return p
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
	log.Panic("not implemented")
	return nil
}

func (f *polyField) Random() Element {
	log.Panic("not implemented")
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
	l := len(p.cs)
	if len(p2.cs) > l {
		l = len(p2.cs)
	}
	zero := p.f.Zero()
	for i := l; i >= 0; i-- {
		c1 := zero
		c2 := zero
		if l < len(p.cs) {
			c1 = p.cs[i]
		}
		if l < len(p2.cs) {
			c2 = p2.cs[i]
		}
		cmp := c1.Cmp(c2)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}
