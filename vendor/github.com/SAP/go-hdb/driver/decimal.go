/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
)

//bigint word size (*--> src/pkg/math/big/arith.go)
const (
	// Compute the size _S of a Word in bytes.
	_m    = ^big.Word(0)
	_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S    = 1 << _logS
)

const (
	// http://en.wikipedia.org/wiki/Decimal128_floating-point_format
	dec128Digits = 34
	dec128Bias   = 6176
	dec128MinExp = -6176
	dec128MaxExp = 6111
)

const (
	decimalSize = 16 //number of bytes
)

var natZero = big.NewInt(0)
var natOne = big.NewInt(1)
var natTen = big.NewInt(10)

var nat = []*big.Int{
	natOne,                  //10^0
	natTen,                  //10^1
	big.NewInt(100),         //10^2
	big.NewInt(1000),        //10^3
	big.NewInt(10000),       //10^4
	big.NewInt(100000),      //10^5
	big.NewInt(1000000),     //10^6
	big.NewInt(10000000),    //10^7
	big.NewInt(100000000),   //10^8
	big.NewInt(1000000000),  //10^9
	big.NewInt(10000000000), //10^10
}

const lg10 = math.Ln10 / math.Ln2 // ~log2(10)

var maxDecimal = new(big.Int).SetBytes([]byte{0x01, 0xED, 0x09, 0xBE, 0xAD, 0x87, 0xC0, 0x37, 0x8D, 0x8E, 0x63, 0xFF, 0xFF, 0xFF, 0xFF})

type decFlags byte

const (
	dfNotExact decFlags = 1 << iota
	dfOverflow
	dfUnderflow
)

// ErrDecimalOutOfRange means that a big.Rat exceeds the size of hdb decimal fields.
var ErrDecimalOutOfRange = errors.New("decimal out of range error")

// big.Int free list
var bigIntFree = sync.Pool{
	New: func() interface{} { return new(big.Int) },
}

// big.Rat free list
var bigRatFree = sync.Pool{
	New: func() interface{} { return new(big.Rat) },
}

// A Decimal is the driver representation of a database decimal field value as big.Rat.
type Decimal big.Rat

// Scan implements the database/sql/Scanner interface.
func (d *Decimal) Scan(src interface{}) error {

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("decimal: invalid data type %T", src)
	}

	if len(b) != decimalSize {
		return fmt.Errorf("decimal: invalid size %d of %v - %d expected", len(b), b, decimalSize)
	}

	if (b[15] & 0x60) == 0x60 {
		return fmt.Errorf("decimal: format (infinity, nan, ...) not supported : %v", b)
	}

	v := (*big.Rat)(d)
	p := v.Num()
	q := v.Denom()

	neg, exp := decodeDecimal(b, p)

	switch {
	case exp < 0:
		q.Set(exp10(exp * -1))
	case exp == 0:
		q.Set(natOne)
	case exp > 0:
		p.Mul(p, exp10(exp))
		q.Set(natOne)
	}

	if neg {
		v.Neg(v)
	}
	return nil
}

// Value implements the database/sql/Valuer interface.
func (d Decimal) Value() (driver.Value, error) {
	m := bigIntFree.Get().(*big.Int)
	neg, exp, df := convertRatToDecimal((*big.Rat)(&d), m, dec128Digits, dec128MinExp, dec128MaxExp)

	var v driver.Value
	var err error

	switch {
	default:
		v, err = encodeDecimal(m, neg, exp)
	case df&dfUnderflow != 0: // set to zero
		m.Set(natZero)
		v, err = encodeDecimal(m, false, 0)
	case df&dfOverflow != 0:
		err = ErrDecimalOutOfRange
	}

	// performance (avoid expensive defer)
	bigIntFree.Put(m)

	return v, err
}

func convertRatToDecimal(x *big.Rat, m *big.Int, digits, minExp, maxExp int) (bool, int, decFlags) {

	neg := x.Sign() < 0 //store sign

	if x.Num().Cmp(natZero) == 0 { // zero
		m.Set(natZero)
		return neg, 0, 0
	}

	c := bigRatFree.Get().(*big.Rat).Abs(x) // copy && abs
	a := c.Num()
	b := c.Denom()

	exp, shift := 0, 0

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

	var df decFlags

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

	// performance (avoid expensive defer)
	bigRatFree.Put(c)

	return neg, exp, df
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// performance: tested with reference work variable
// - but int.Set is expensive, so let's live with big.Int creation for n >= len(nat)
func exp10(n int) *big.Int {
	if n < len(nat) {
		return nat[n]
	}
	r := big.NewInt(int64(n))
	return r.Exp(natTen, r, nil)
}

func digits10(p *big.Int) int {
	k := p.BitLen() // 2^k <= p < 2^(k+1) - 1
	//i := int(float64(k) / lg10) //minimal digits base 10
	//i := int(float64(k) / lg10) //minimal digits base 10
	i := k * 100 / 332
	if i < 1 {
		i = 1
	}

	for ; ; i++ {
		if p.Cmp(exp10(i)) < 0 {
			return i
		}
	}
}

func decodeDecimal(b []byte, m *big.Int) (bool, int) {

	neg := (b[15] & 0x80) != 0
	exp := int((((uint16(b[15])<<8)|uint16(b[14]))<<1)>>2) - dec128Bias

	b14 := b[14]  // save b[14]
	b[14] &= 0x01 // keep the mantissa bit (rest: sign and exp)

	//most significand byte
	msb := 14
	for msb > 0 {
		if b[msb] != 0 {
			break
		}
		msb--
	}

	//calc number of words
	numWords := (msb / _S) + 1
	w := make([]big.Word, numWords)

	k := numWords - 1
	d := big.Word(0)
	for i := msb; i >= 0; i-- {
		d |= big.Word(b[i])
		if k*_S == i {
			w[k] = d
			k--
			d = 0
		}
		d <<= 8
	}
	b[14] = b14 // restore b[14]
	m.SetBits(w)
	return neg, exp
}

func encodeDecimal(m *big.Int, neg bool, exp int) (driver.Value, error) {

	b := make([]byte, decimalSize)

	// little endian bigint words (significand) -> little endian db decimal format
	j := 0
	for _, d := range m.Bits() {
		for i := 0; i < 8; i++ {
			b[j] = byte(d)
			d >>= 8
			j++
		}
	}

	exp += dec128Bias
	b[14] |= (byte(exp) << 1)
	b[15] = byte(uint16(exp) >> 7)

	if neg {
		b[15] |= 0x80
	}

	return b, nil
}

// NullDecimal represents an Decimal that may be null.
// NullDecimal implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullDecimal struct {
	Decimal *Decimal
	Valid   bool // Valid is true if Decimal is not NULL
}

// Scan implements the Scanner interface.
func (n *NullDecimal) Scan(value interface{}) error {
	var b []byte

	b, n.Valid = value.([]byte)
	if !n.Valid {
		return nil
	}
	if n.Decimal == nil {
		return fmt.Errorf("invalid decimal value %v", n.Decimal)
	}
	return n.Decimal.Scan(b)
}

// Value implements the driver Valuer interface.
func (n NullDecimal) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	if n.Decimal == nil {
		return nil, fmt.Errorf("invalid decimal value %v", n.Decimal)
	}
	return n.Decimal.Value()
}
