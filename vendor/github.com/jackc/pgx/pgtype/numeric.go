package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

// PostgreSQL internal numeric storage uses 16-bit "digits" with base of 10,000
const nbase = 10000

var big0 *big.Int = big.NewInt(0)
var big1 *big.Int = big.NewInt(1)
var big10 *big.Int = big.NewInt(10)
var big100 *big.Int = big.NewInt(100)
var big1000 *big.Int = big.NewInt(1000)

var bigMaxInt8 *big.Int = big.NewInt(math.MaxInt8)
var bigMinInt8 *big.Int = big.NewInt(math.MinInt8)
var bigMaxInt16 *big.Int = big.NewInt(math.MaxInt16)
var bigMinInt16 *big.Int = big.NewInt(math.MinInt16)
var bigMaxInt32 *big.Int = big.NewInt(math.MaxInt32)
var bigMinInt32 *big.Int = big.NewInt(math.MinInt32)
var bigMaxInt64 *big.Int = big.NewInt(math.MaxInt64)
var bigMinInt64 *big.Int = big.NewInt(math.MinInt64)
var bigMaxInt *big.Int = big.NewInt(int64(maxInt))
var bigMinInt *big.Int = big.NewInt(int64(minInt))

var bigMaxUint8 *big.Int = big.NewInt(math.MaxUint8)
var bigMaxUint16 *big.Int = big.NewInt(math.MaxUint16)
var bigMaxUint32 *big.Int = big.NewInt(math.MaxUint32)
var bigMaxUint64 *big.Int = (&big.Int{}).SetUint64(uint64(math.MaxUint64))
var bigMaxUint *big.Int = (&big.Int{}).SetUint64(uint64(maxUint))

var bigNBase *big.Int = big.NewInt(nbase)
var bigNBaseX2 *big.Int = big.NewInt(nbase * nbase)
var bigNBaseX3 *big.Int = big.NewInt(nbase * nbase * nbase)
var bigNBaseX4 *big.Int = big.NewInt(nbase * nbase * nbase * nbase)

type Numeric struct {
	Int    *big.Int
	Exp    int32
	Status Status
}

func (dst *Numeric) Set(src interface{}) error {
	if src == nil {
		*dst = Numeric{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case float32:
		num, exp, err := parseNumericString(strconv.FormatFloat(float64(value), 'f', -1, 64))
		if err != nil {
			return err
		}
		*dst = Numeric{Int: num, Exp: exp, Status: Present}
	case float64:
		num, exp, err := parseNumericString(strconv.FormatFloat(value, 'f', -1, 64))
		if err != nil {
			return err
		}
		*dst = Numeric{Int: num, Exp: exp, Status: Present}
	case int8:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case uint8:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case int16:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case uint16:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case int32:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case uint32:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case int64:
		*dst = Numeric{Int: big.NewInt(value), Status: Present}
	case uint64:
		*dst = Numeric{Int: (&big.Int{}).SetUint64(value), Status: Present}
	case int:
		*dst = Numeric{Int: big.NewInt(int64(value)), Status: Present}
	case uint:
		*dst = Numeric{Int: (&big.Int{}).SetUint64(uint64(value)), Status: Present}
	case string:
		num, exp, err := parseNumericString(value)
		if err != nil {
			return err
		}
		*dst = Numeric{Int: num, Exp: exp, Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Numeric", value)
	}

	return nil
}

func (dst *Numeric) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Numeric) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *float32:
			f, err := src.toFloat64()
			if err != nil {
				return err
			}
			return float64AssignTo(f, src.Status, dst)
		case *float64:
			f, err := src.toFloat64()
			if err != nil {
				return err
			}
			return float64AssignTo(f, src.Status, dst)
		case *int:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(bigMaxInt) > 0 {
				return errors.Errorf("%v is greater than maximum value for %T", normalizedInt, *v)
			}
			if normalizedInt.Cmp(bigMinInt) < 0 {
				return errors.Errorf("%v is less than minimum value for %T", normalizedInt, *v)
			}
			*v = int(normalizedInt.Int64())
		case *int8:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(bigMaxInt8) > 0 {
				return errors.Errorf("%v is greater than maximum value for %T", normalizedInt, *v)
			}
			if normalizedInt.Cmp(bigMinInt8) < 0 {
				return errors.Errorf("%v is less than minimum value for %T", normalizedInt, *v)
			}
			*v = int8(normalizedInt.Int64())
		case *int16:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(bigMaxInt16) > 0 {
				return errors.Errorf("%v is greater than maximum value for %T", normalizedInt, *v)
			}
			if normalizedInt.Cmp(bigMinInt16) < 0 {
				return errors.Errorf("%v is less than minimum value for %T", normalizedInt, *v)
			}
			*v = int16(normalizedInt.Int64())
		case *int32:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(bigMaxInt32) > 0 {
				return errors.Errorf("%v is greater than maximum value for %T", normalizedInt, *v)
			}
			if normalizedInt.Cmp(bigMinInt32) < 0 {
				return errors.Errorf("%v is less than minimum value for %T", normalizedInt, *v)
			}
			*v = int32(normalizedInt.Int64())
		case *int64:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(bigMaxInt64) > 0 {
				return errors.Errorf("%v is greater than maximum value for %T", normalizedInt, *v)
			}
			if normalizedInt.Cmp(bigMinInt64) < 0 {
				return errors.Errorf("%v is less than minimum value for %T", normalizedInt, *v)
			}
			*v = normalizedInt.Int64()
		case *uint:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(big0) < 0 {
				return errors.Errorf("%d is less than zero for %T", normalizedInt, *v)
			} else if normalizedInt.Cmp(bigMaxUint) > 0 {
				return errors.Errorf("%d is greater than maximum value for %T", normalizedInt, *v)
			}
			*v = uint(normalizedInt.Uint64())
		case *uint8:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(big0) < 0 {
				return errors.Errorf("%d is less than zero for %T", normalizedInt, *v)
			} else if normalizedInt.Cmp(bigMaxUint8) > 0 {
				return errors.Errorf("%d is greater than maximum value for %T", normalizedInt, *v)
			}
			*v = uint8(normalizedInt.Uint64())
		case *uint16:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(big0) < 0 {
				return errors.Errorf("%d is less than zero for %T", normalizedInt, *v)
			} else if normalizedInt.Cmp(bigMaxUint16) > 0 {
				return errors.Errorf("%d is greater than maximum value for %T", normalizedInt, *v)
			}
			*v = uint16(normalizedInt.Uint64())
		case *uint32:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(big0) < 0 {
				return errors.Errorf("%d is less than zero for %T", normalizedInt, *v)
			} else if normalizedInt.Cmp(bigMaxUint32) > 0 {
				return errors.Errorf("%d is greater than maximum value for %T", normalizedInt, *v)
			}
			*v = uint32(normalizedInt.Uint64())
		case *uint64:
			normalizedInt, err := src.toBigInt()
			if err != nil {
				return err
			}
			if normalizedInt.Cmp(big0) < 0 {
				return errors.Errorf("%d is less than zero for %T", normalizedInt, *v)
			} else if normalizedInt.Cmp(bigMaxUint64) > 0 {
				return errors.Errorf("%d is greater than maximum value for %T", normalizedInt, *v)
			}
			*v = normalizedInt.Uint64()
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return nil
}

func (dst *Numeric) toBigInt() (*big.Int, error) {
	if dst.Exp == 0 {
		return dst.Int, nil
	}

	num := &big.Int{}
	num.Set(dst.Int)
	if dst.Exp > 0 {
		mul := &big.Int{}
		mul.Exp(big10, big.NewInt(int64(dst.Exp)), nil)
		num.Mul(num, mul)
		return num, nil
	}

	div := &big.Int{}
	div.Exp(big10, big.NewInt(int64(-dst.Exp)), nil)
	remainder := &big.Int{}
	num.DivMod(num, div, remainder)
	if remainder.Cmp(big0) != 0 {
		return nil, errors.Errorf("cannot convert %v to integer", dst)
	}
	return num, nil
}

func (src *Numeric) toFloat64() (float64, error) {
	f, err := strconv.ParseFloat(src.Int.String(), 64)
	if err != nil {
		return 0, err
	}
	if src.Exp > 0 {
		for i := 0; i < int(src.Exp); i++ {
			f *= 10
		}
	} else if src.Exp < 0 {
		for i := 0; i > int(src.Exp); i-- {
			f /= 10
		}
	}
	return f, nil
}

func (dst *Numeric) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Numeric{Status: Null}
		return nil
	}

	num, exp, err := parseNumericString(string(src))
	if err != nil {
		return err
	}

	*dst = Numeric{Int: num, Exp: exp, Status: Present}
	return nil
}

func parseNumericString(str string) (n *big.Int, exp int32, err error) {
	parts := strings.SplitN(str, ".", 2)
	digits := strings.Join(parts, "")

	if len(parts) > 1 {
		exp = int32(-len(parts[1]))
	} else {
		for len(digits) > 1 && digits[len(digits)-1] == '0' {
			digits = digits[:len(digits)-1]
			exp++
		}
	}

	accum := &big.Int{}
	if _, ok := accum.SetString(digits, 10); !ok {
		return nil, 0, errors.Errorf("%s is not a number", str)
	}

	return accum, exp, nil
}

func (dst *Numeric) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Numeric{Status: Null}
		return nil
	}

	if len(src) < 8 {
		return errors.Errorf("numeric incomplete %v", src)
	}

	rp := 0
	ndigits := int16(binary.BigEndian.Uint16(src[rp:]))
	rp += 2

	if ndigits == 0 {
		*dst = Numeric{Int: big.NewInt(0), Status: Present}
		return nil
	}

	weight := int16(binary.BigEndian.Uint16(src[rp:]))
	rp += 2
	sign := int16(binary.BigEndian.Uint16(src[rp:]))
	rp += 2
	dscale := int16(binary.BigEndian.Uint16(src[rp:]))
	rp += 2

	if len(src[rp:]) < int(ndigits)*2 {
		return errors.Errorf("numeric incomplete %v", src)
	}

	accum := &big.Int{}

	for i := 0; i < int(ndigits+3)/4; i++ {
		int64accum, bytesRead, digitsRead := nbaseDigitsToInt64(src[rp:])
		rp += bytesRead

		if i > 0 {
			var mul *big.Int
			switch digitsRead {
			case 1:
				mul = bigNBase
			case 2:
				mul = bigNBaseX2
			case 3:
				mul = bigNBaseX3
			case 4:
				mul = bigNBaseX4
			default:
				return errors.Errorf("invalid digitsRead: %d (this can't happen)", digitsRead)
			}
			accum.Mul(accum, mul)
		}

		accum.Add(accum, big.NewInt(int64accum))
	}

	exp := (int32(weight) - int32(ndigits) + 1) * 4

	if dscale > 0 {
		fracNBaseDigits := ndigits - weight - 1
		fracDecimalDigits := fracNBaseDigits * 4

		if dscale > fracDecimalDigits {
			multCount := int(dscale - fracDecimalDigits)
			for i := 0; i < multCount; i++ {
				accum.Mul(accum, big10)
				exp--
			}
		} else if dscale < fracDecimalDigits {
			divCount := int(fracDecimalDigits - dscale)
			for i := 0; i < divCount; i++ {
				accum.Div(accum, big10)
				exp++
			}
		}
	}

	reduced := &big.Int{}
	remainder := &big.Int{}
	if exp >= 0 {
		for {
			reduced.DivMod(accum, big10, remainder)
			if remainder.Cmp(big0) != 0 {
				break
			}
			accum.Set(reduced)
			exp++
		}
	}

	if sign != 0 {
		accum.Neg(accum)
	}

	*dst = Numeric{Int: accum, Exp: exp, Status: Present}

	return nil

}

func nbaseDigitsToInt64(src []byte) (accum int64, bytesRead, digitsRead int) {
	digits := len(src) / 2
	if digits > 4 {
		digits = 4
	}

	rp := 0

	for i := 0; i < digits; i++ {
		if i > 0 {
			accum *= nbase
		}
		accum += int64(binary.BigEndian.Uint16(src[rp:]))
		rp += 2
	}

	return accum, rp, digits
}

func (src *Numeric) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, src.Int.String()...)
	buf = append(buf, 'e')
	buf = append(buf, strconv.FormatInt(int64(src.Exp), 10)...)
	return buf, nil
}

func (src *Numeric) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var sign int16
	if src.Int.Cmp(big0) < 0 {
		sign = 16384
	}

	absInt := &big.Int{}
	wholePart := &big.Int{}
	fracPart := &big.Int{}
	remainder := &big.Int{}
	absInt.Abs(src.Int)

	// Normalize absInt and exp to where exp is always a multiple of 4. This makes
	// converting to 16-bit base 10,000 digits easier.
	var exp int32
	switch src.Exp % 4 {
	case 1, -3:
		exp = src.Exp - 1
		absInt.Mul(absInt, big10)
	case 2, -2:
		exp = src.Exp - 2
		absInt.Mul(absInt, big100)
	case 3, -1:
		exp = src.Exp - 3
		absInt.Mul(absInt, big1000)
	default:
		exp = src.Exp
	}

	if exp < 0 {
		divisor := &big.Int{}
		divisor.Exp(big10, big.NewInt(int64(-exp)), nil)
		wholePart.DivMod(absInt, divisor, fracPart)
		fracPart.Add(fracPart, divisor)
	} else {
		wholePart = absInt
	}

	var wholeDigits, fracDigits []int16

	for wholePart.Cmp(big0) != 0 {
		wholePart.DivMod(wholePart, bigNBase, remainder)
		wholeDigits = append(wholeDigits, int16(remainder.Int64()))
	}

	if fracPart.Cmp(big0) != 0 {
		for fracPart.Cmp(big1) != 0 {
			fracPart.DivMod(fracPart, bigNBase, remainder)
			fracDigits = append(fracDigits, int16(remainder.Int64()))
		}
	}

	buf = pgio.AppendInt16(buf, int16(len(wholeDigits)+len(fracDigits)))

	var weight int16
	if len(wholeDigits) > 0 {
		weight = int16(len(wholeDigits) - 1)
		if exp > 0 {
			weight += int16(exp / 4)
		}
	} else {
		weight = int16(exp/4) - 1 + int16(len(fracDigits))
	}
	buf = pgio.AppendInt16(buf, weight)

	buf = pgio.AppendInt16(buf, sign)

	var dscale int16
	if src.Exp < 0 {
		dscale = int16(-src.Exp)
	}
	buf = pgio.AppendInt16(buf, dscale)

	for i := len(wholeDigits) - 1; i >= 0; i-- {
		buf = pgio.AppendInt16(buf, wholeDigits[i])
	}

	for i := len(fracDigits) - 1; i >= 0; i-- {
		buf = pgio.AppendInt16(buf, fracDigits[i])
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Numeric) Scan(src interface{}) error {
	if src == nil {
		*dst = Numeric{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case float64:
		// TODO
		// *dst = Numeric{Float: src, Status: Present}
		return nil
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Numeric) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		buf, err := src.EncodeText(nil, nil)
		if err != nil {
			return nil, err
		}

		return string(buf), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
