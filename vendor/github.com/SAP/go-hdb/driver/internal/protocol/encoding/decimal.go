package encoding

import (
	"errors"
	"math"
	"math/big"
	"math/bits"
)

// ErrDecimalOutOfRange means that a big.Rat exceeds the size of hdb decimal fields.
var ErrDecimalOutOfRange = errors.New("decimal out of range error")

const _S = bits.UintSize / 8 // word size in bytes
// http://en.wikipedia.org/wiki/Decimal128_floating-point_format
const dec128Bias = 6176
const decSize = 16

// decimals.
const (
	// http://en.wikipedia.org/wiki/Decimal128_floating-point_format
	dec128Digits = 34
	// 	dec128Bias   = 6176
	dec128MinExp = -6176
	dec128MaxExp = 6111
)

var (
	natZero = big.NewInt(0)
	natOne  = big.NewInt(1)
	natTen  = big.NewInt(10)
)

const maxNatExp10 = 38 // maximal fixed decimal precision

var natExp10 = make([]*big.Int, maxNatExp10)

func init() {
	natExp10[0], natExp10[1] = natOne, natTen
	for i := 2; i < maxNatExp10; i++ {
		natExp10[i] = new(big.Int).Mul(natExp10[i-1], natTen)
	}
}

/*
performance: tested with reference work variable
  - but int.Set is expensive, so let's live with big.Int creation for n >= len(nat)
*/
func exp10(n int) *big.Int {
	if n < len(natExp10) {
		return natExp10[n]
	}
	r := big.NewInt(int64(n))
	return r.Exp(natTen, r, nil)
}

var lg10 = math.Log2(10)

func digits10(p *big.Int) int {
	k := p.BitLen() // 2^k <= p < 2^(k+1) - 1
	i := int(float64(k) / lg10)
	if i < 1 {
		i = 1
	}
	// i <= digit10(p)
	for ; ; i++ {
		if p.Cmp(exp10(i)) < 0 {
			return i
		}
	}
}

// decimal flag.
const (
	dfNotExact byte = 1 << iota
	dfOverflow
	dfUnderflow
)

func convertRatToDecimal(x *big.Rat, m *big.Int, digits, minExp, maxExp int) (int, byte) {
	if x.Num().Cmp(natZero) == 0 { // zero
		m.Set(natZero)
		return 0, 0
	}

	var tmp big.Rat

	c := (&tmp).Set(x) // copy
	a := c.Num()
	b := c.Denom()

	var exp int
	shift := 0

	if c.IsInt() {
		exp = digits10(a) - 1
	} else {
		shift = digits10(a) - digits10(b)
		switch {
		case shift < 0:
			a.Mul(a, exp10(shift*-1))
		case shift > 0:
			b.Mul(b, exp10(shift))
		}
		if a.Cmp(b) == -1 {
			exp = shift - 1
		} else {
			exp = shift
		}
	}

	var df byte

	switch {
	default:
		exp = max(exp-digits+1, minExp)
	case exp < minExp:
		df |= dfUnderflow
		exp = exp - digits + 1
	}

	if exp > maxExp {
		df |= dfOverflow
	}

	shift = exp - shift
	switch {
	case shift < 0:
		a.Mul(a, exp10(shift*-1))
	case exp > 0:
		b.Mul(b, exp10(shift))
	}

	m.QuoRem(a, b, a) // reuse a as rest
	if a.Cmp(natZero) != 0 {
		// round (business >= 0.5 up)
		df |= dfNotExact
		if a.Add(a, a).Cmp(b) >= 0 {
			m.Add(m, natOne)
			if m.Cmp(exp10(digits)) == 0 {
				shift := min(digits, maxExp-exp)
				if shift < 1 { // overflow -> shift one at minimum
					df |= dfOverflow
					shift = 1
				}
				m.Set(exp10(digits - shift))
				exp += shift
			}
		}
	}

	// norm
	for exp < maxExp {
		a.QuoRem(m, natTen, b) // reuse a, b
		if b.Cmp(natZero) != 0 {
			break
		}
		m.Set(a)
		exp++
	}

	return exp, df
}

func convertDecimalToRat(m *big.Int, exp int) *big.Rat {
	if m == nil {
		return nil
	}

	v := new(big.Rat).SetInt(m)
	p := v.Num()
	q := v.Denom()

	switch {
	case exp < 0:
		q.Set(exp10(exp * -1))
	case exp == 0:
		q.Set(natOne)
	case exp > 0:
		p.Mul(p, exp10(exp))
		q.Set(natOne)
	}
	return v
}

func convertRatToFixed(r *big.Rat, m *big.Int, prec, scale int) byte {
	if scale < 0 {
		panic("fixed: invalid scale")
	}

	var df byte

	m.Set(r.Num())
	m.Mul(m, exp10(scale))

	var tmp big.Rat

	c := (&tmp).SetFrac(m, r.Denom()) // norm
	a := c.Num()
	b := c.Denom()

	if b.Cmp(natZero) == 0 { //
		m.Set(a)
		return df
	}

	m.QuoRem(a, b, a) // reuse a as rest
	if a.Cmp(natZero) != 0 {
		// round (business >= 0.5 up)
		df |= dfNotExact
		if a.Add(a, a).Cmp(b) >= 0 {
			m.Add(m, natOne)
		}
	}

	maxInt := exp10(prec)
	minInt := new(big.Int).Neg(maxInt)

	if m.Cmp(minInt) <= 0 || m.Cmp(maxInt) >= 0 {
		df |= dfOverflow
	}
	return df
}

func convertFixedToRat(m *big.Int, scale int) *big.Rat {
	if m == nil {
		return nil
	}
	if scale < 0 {
		panic("fixed: invalid scale")
	}
	q := exp10(scale)
	return new(big.Rat).SetFrac(m, q)
}
