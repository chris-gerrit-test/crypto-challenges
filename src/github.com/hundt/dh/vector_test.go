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
