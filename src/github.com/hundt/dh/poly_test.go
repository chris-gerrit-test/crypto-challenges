package dh

import (
	"math/big"
	"testing"
)

func TestPolyFieldAdd(t *testing.T) {
	g := NewFiniteGroup(*big.NewInt(7))

	p := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(5)),
		NewFiniteElement(g, *big.NewInt(2)),
	})

	p1 := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(4)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(6)),
	})

	p2 := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(1)),
		NewFiniteElement(g, *big.NewInt(2)),
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(4)),
	})

	f := NewPolyField(p)

	p3 := f.Add(p1, p2)

	if p3.String() != "3*x^3 + 3*x^2 + 6*x + 4" {
		t.Errorf("%s + %s = %s", p1, p2, p3)
	}

	if f.Sub(p3, p1).Cmp(p2) != 0 {
		t.Errorf("%s - (%s) = %s", p3, p1, f.Sub(p3, p1))
	}

	if f.Sub(p3, p2).Cmp(p1) != 0 {
		t.Errorf("%s - (%s) = %s", p3, p2, f.Sub(p3, p2))
	}

	g = NewFiniteGroup(*big.NewInt(251))

	p = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(7)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(9)),
		NewFiniteElement(g, *big.NewInt(12)),
		NewFiniteElement(g, *big.NewInt(1)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p1 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(4)),
		NewFiniteElement(g, *big.NewInt(7)),
		NewFiniteElement(g, *big.NewInt(76)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(123)),
	})

	p2 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(76)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(225)),
		NewFiniteElement(g, *big.NewInt(12)),
		NewFiniteElement(g, *big.NewInt(196)),
	})

	f = NewPolyField(p)

	p3 = f.Add(p1, p2)

	if p3.String() != "68*x^4 + 12*x^3 + 50*x^2 + 7*x + 80" {
		t.Errorf("%s + %s = %s", p1, p2, p3)
	}

	if f.Sub(p3, p1).Cmp(p2) != 0 {
		t.Errorf("%s - (%s) = %s", p3, p1, f.Sub(p3, p1))
	}

	if f.Sub(p3, p2).Cmp(p1) != 0 {
		t.Errorf("%s - (%s) = %s", p3, p2, f.Sub(p3, p2))
	}
}

func TestPolyFieldMul(t *testing.T) {
	g := NewFiniteGroup(*big.NewInt(7))

	p := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p1 := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p2 := NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(4)),
		NewFiniteElement(g, *big.NewInt(2)),
	})

	f := NewPolyField(p)

	p3 := f.Op(p1, p2)

	if p3.String() != "3*x + 6" {
		t.Errorf("(%s) * (%s) mod (%s) = %s", p1, p2, p, p3)
	}

	p = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(3)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(2)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p1 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(4)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p2 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(1)),
		NewFiniteElement(g, *big.NewInt(1)),
		NewFiniteElement(g, *big.NewInt(0)),
	})

	f = NewPolyField(p)

	p3 = f.Op(p1, p2)

	if p3.String() != "6*x^2 + 4*x + 1" {
		t.Errorf("(%s) * (%s) mod (%s) = %s", p1, p2, p, p3)
	}

	g = NewFiniteGroup(*big.NewInt(251))

	p = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(7)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(9)),
		NewFiniteElement(g, *big.NewInt(12)),
		NewFiniteElement(g, *big.NewInt(1)),
		NewFiniteElement(g, *big.NewInt(1)),
	})

	p1 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(4)),
		NewFiniteElement(g, *big.NewInt(7)),
		NewFiniteElement(g, *big.NewInt(76)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(123)),
	})

	p2 = NewPolyElement(g, []Element{
		NewFiniteElement(g, *big.NewInt(76)),
		NewFiniteElement(g, *big.NewInt(0)),
		NewFiniteElement(g, *big.NewInt(225)),
		NewFiniteElement(g, *big.NewInt(12)),
		NewFiniteElement(g, *big.NewInt(196)),
	})

	f = NewPolyField(p)

	p3 = f.Op(p1, p2)

	if p3.String() != "117*x^4 + 151*x^3 + 117*x^2 + 182*x + 217" {
		t.Errorf("%s * (%s) = %s", p1, p2, p3)
	}
}
