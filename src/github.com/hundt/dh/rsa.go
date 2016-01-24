package dh

import "math/big"

type rsa struct {
	d, e, n *big.Int
}

type RSA interface {
	Encrypt(m []byte) *big.Int
	Decrypt(c *big.Int) []byte
	PublicKey() *big.Int
}

func NewRSA(p, q *big.Int) RSA {
	n := new(big.Int).Mul(p, q)
	t := new(big.Int).Sub(n, p)
	t.Sub(t, q)
	t.Add(t, big.NewInt(1))
	e := big.NewInt(65537)
	d := new(big.Int).ModInverse(e, t)
	return &rsa{d, e, n}
}

func (r *rsa) Encrypt(m []byte) *big.Int {
	z := new(big.Int).SetBytes(m)
	return z.Exp(z, r.e, r.n)
}

func (r *rsa) Decrypt(c *big.Int) []byte {
	return new(big.Int).Exp(c, r.d, r.n).Bytes()
}

func (r *rsa) PublicKey() *big.Int {
	return new(big.Int).Set(r.e)
}
