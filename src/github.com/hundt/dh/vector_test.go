package dh

import (
	"math/big"
	"testing"
)

func TestEq(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(3, 4)})
	v2 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(3, 5)})
	if v1.Eq(v2) {
		t.Errorf("%s should not equal %s", v1, v2)
	}
	if v2.Eq(v1) {
		t.Errorf("%s should not equal %s", v1, v2)
	}
	v2 = Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(3, 4)})
	if !v1.Eq(v2) {
		t.Errorf("%s should equal %s", v1, v2)
	}
	if !v2.Eq(v1) {
		t.Errorf("%s should equal %s", v1, v2)
	}
	v1 = Vector([]*big.Rat{
		big.NewRat(1, 3),
		big.NewRat(3, 4)})
	if v1.Eq(v2) {
		t.Errorf("%s should not equal %s", v1, v2)
	}
	if v2.Eq(v1) {
		t.Errorf("%s should not equal %s", v1, v2)
	}
}

func TestDot(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(2, 3)})
	v2 := Vector([]*big.Rat{
		big.NewRat(1, 3),
		big.NewRat(5, 4)})
	actual := v1.Dot(v2)
	expected := big.NewRat(1, 1)
	if expected.Cmp(actual) != 0 {
		t.Errorf("%s.%s should equal %s, not %s", v1, v2, expected, actual)
	}
}

func TestFastDot(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(2, 3)})
	v2 := Vector([]*big.Rat{
		big.NewRat(1, 3),
		big.NewRat(5, 4)})
	actual := v1.fastDot(v2)
	expected := big.NewRat(1, 1)
	if expected.Cmp(actual) != 0 {
		t.Errorf("%s.%s should equal %s, not %s", v1, v2, expected, actual)
	}
}

func TestScale(t *testing.T) {
	v := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(2, 3)})
	s := big.NewRat(3, 4)
	actual := v.Scale(s)
	expected := Vector([]*big.Rat{
		big.NewRat(3, 8),
		big.NewRat(1, 2)})
	if !expected.Eq(actual) {
		t.Errorf("%s%s should equal %s, not %s", s, v, expected, actual)
	}
}

func TestAdd(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(3, 4)})
	v2 := Vector([]*big.Rat{
		big.NewRat(1, 3),
		big.NewRat(3, 5)})

	v3 := v1.Copy()
	v4 := Vector([]*big.Rat{
		big.NewRat(5, 6),
		big.NewRat(27, 20)})
	v3.Add(v2)
	if !v3.Eq(v4) {
		t.Errorf("%s + %s should equal %s, not %s", v1, v2, v4, v3)
	}

	v3 = v2.Copy()
	v3.Add(v1)
	if !v3.Eq(v4) {
		t.Errorf("%s + %s should equal %s, not %s", v2, v1, v4, v3)
	}
}

func TestSub(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(1, 2),
		big.NewRat(3, 4)})
	v2 := Vector([]*big.Rat{
		big.NewRat(1, 3),
		big.NewRat(3, 5)})

	v3 := v1.Copy()
	v4 := Vector([]*big.Rat{
		big.NewRat(1, 6),
		big.NewRat(3, 20)})
	v3.Sub(v2)
	if !v3.Eq(v4) {
		t.Errorf("%s - %s should equal %s, not %s", v1, v2, v4, v3)
	}

	v3 = v2.Copy()
	v4 = Vector([]*big.Rat{
		big.NewRat(-1, 6),
		big.NewRat(-3, 20)})
	v3.Sub(v1)
	if !v3.Eq(v4) {
		t.Errorf("%s - %s should equal %s, not %s", v2, v1, v4, v3)
	}
}

func TestProject(t *testing.T) {
	v1 := Vector([]*big.Rat{
		big.NewRat(2, 1),
		big.NewRat(1, 1)})
	v2 := Vector([]*big.Rat{
		big.NewRat(-3, 1),
		big.NewRat(4, 1)})

	v3 := v1.Copy()
	v4 := Vector([]*big.Rat{
		big.NewRat(6, 25),
		big.NewRat(-8, 25)})
	v3.Proj(v2)
	if !v3.Eq(v4) {
		t.Errorf("%s projected onto %s should equal %s, not %s", v1, v2, v4, v3)
	}

	v1 = Vector([]*big.Rat{
		big.NewRat(-2, 1),
		big.NewRat(3, 1),
		big.NewRat(12, 1)})
	v2 = Vector([]*big.Rat{
		big.NewRat(-1, 2),
		big.NewRat(1, 3),
		big.NewRat(4, 5)})

	v3 = v1.Copy()
	v4 = Vector([]*big.Rat{
		big.NewRat(-5220, 901),
		big.NewRat(3480, 901),
		big.NewRat(8352, 901)})
	v3.Proj(v2)
	if !v3.Eq(v4) {
		t.Errorf("%s projected onto %s should equal %s, not %s", v1, v2, v4, v3)
	}
}

func TestRound(t *testing.T) {
	r1 := big.NewRat(1, 2)
	actual := Round(new(big.Rat).Set(r1))
	expected := big.NewRat(1, 1)
	if actual.Cmp(expected) != 0 {
		t.Errorf("round(%s) should equal %s, not %s", r1, expected, actual)
	}
	r1 = big.NewRat(1, 3)
	actual = Round(new(big.Rat).Set(r1))
	expected = big.NewRat(0, 1)
	if actual.Cmp(expected) != 0 {
		t.Errorf("round(%s) should equal %s, not %s", r1, expected, actual)
	}
	r1 = big.NewRat(5, 3)
	actual = Round(new(big.Rat).Set(r1))
	expected = big.NewRat(2, 1)
	if actual.Cmp(expected) != 0 {
		t.Errorf("round(%s) should equal %s, not %s", r1, expected, actual)
	}
}

func TestLLL(t *testing.T) {
	// Cryptopals example
	m := Matrix([]Vector{
		Vector([]*big.Rat{
			big.NewRat(-2, 1),
			big.NewRat(0, 1),
			big.NewRat(2, 1),
			big.NewRat(0, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(1, 2),
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(0, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(-2, 1),
			big.NewRat(1, 2),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 1),
			big.NewRat(1, 1),
			big.NewRat(1, 1),
			big.NewRat(2, 1),
		}),
	})
	expected := Matrix([]Vector{
		Vector([]*big.Rat{
			big.NewRat(1, 2),
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(0, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(-2, 1),
			big.NewRat(1, 2),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 2),
			big.NewRat(0, 1),
			big.NewRat(1, 1),
			big.NewRat(2, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(-3, 2),
			big.NewRat(-1, 1),
			big.NewRat(2, 1),
			big.NewRat(0, 1),
		}),
	})
	actual := m.Copy().LLL(big.NewRat(99, 100))
	if !actual.Eq(expected) {
		t.Errorf("LLL of:\n%s should equal:\n%s not:\n%s", m, expected, actual)
	}

	// Wikipedia example
	m = Matrix([]Vector{
		Vector([]*big.Rat{
			big.NewRat(1, 1),
			big.NewRat(1, 1),
			big.NewRat(1, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(2, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(3, 1),
			big.NewRat(5, 1),
			big.NewRat(6, 1),
		}),
	})
	expected = Matrix([]Vector{
		Vector([]*big.Rat{
			big.NewRat(0, 1),
			big.NewRat(1, 1),
			big.NewRat(0, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(1, 1),
			big.NewRat(0, 1),
			big.NewRat(1, 1),
		}),
		Vector([]*big.Rat{
			big.NewRat(-1, 1),
			big.NewRat(0, 1),
			big.NewRat(2, 1),
		}),
	})
	actual = m.Copy().LLL(big.NewRat(99, 100))
	if !actual.Eq(expected) {
		t.Errorf("LLL of:\n%s should equal:\n%s not:\n%s", m, expected, actual)
	}
}
