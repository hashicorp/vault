package pgtype

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
)

type BoundType byte

const (
	Inclusive = BoundType('i')
	Exclusive = BoundType('e')
	Unbounded = BoundType('U')
	Empty     = BoundType('E')
)

func (bt BoundType) String() string {
	return string(bt)
}

type UntypedTextRange struct {
	Lower     string
	Upper     string
	LowerType BoundType
	UpperType BoundType
}

func ParseUntypedTextRange(src string) (*UntypedTextRange, error) {
	utr := &UntypedTextRange{}
	if src == "empty" {
		utr.LowerType = Empty
		utr.UpperType = Empty
		return utr, nil
	}

	buf := bytes.NewBufferString(src)

	skipWhitespace(buf)

	r, _, err := buf.ReadRune()
	if err != nil {
		return nil, errors.Errorf("invalid lower bound: %v", err)
	}
	switch r {
	case '(':
		utr.LowerType = Exclusive
	case '[':
		utr.LowerType = Inclusive
	default:
		return nil, errors.Errorf("missing lower bound, instead got: %v", string(r))
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, errors.Errorf("invalid lower value: %v", err)
	}
	buf.UnreadRune()

	if r == ',' {
		utr.LowerType = Unbounded
	} else {
		utr.Lower, err = rangeParseValue(buf)
		if err != nil {
			return nil, errors.Errorf("invalid lower value: %v", err)
		}
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, errors.Errorf("missing range separator: %v", err)
	}
	if r != ',' {
		return nil, errors.Errorf("missing range separator: %v", r)
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, errors.Errorf("invalid upper value: %v", err)
	}

	if r == ')' || r == ']' {
		utr.UpperType = Unbounded
	} else {
		buf.UnreadRune()
		utr.Upper, err = rangeParseValue(buf)
		if err != nil {
			return nil, errors.Errorf("invalid upper value: %v", err)
		}

		r, _, err = buf.ReadRune()
		if err != nil {
			return nil, errors.Errorf("missing upper bound: %v", err)
		}
		switch r {
		case ')':
			utr.UpperType = Exclusive
		case ']':
			utr.UpperType = Inclusive
		default:
			return nil, errors.Errorf("missing upper bound, instead got: %v", string(r))
		}
	}

	skipWhitespace(buf)

	if buf.Len() > 0 {
		return nil, errors.Errorf("unexpected trailing data: %v", buf.String())
	}

	return utr, nil
}

func rangeParseValue(buf *bytes.Buffer) (string, error) {
	r, _, err := buf.ReadRune()
	if err != nil {
		return "", err
	}
	if r == '"' {
		return rangeParseQuotedValue(buf)
	}
	buf.UnreadRune()

	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case '\\':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
		case ',', '[', ']', '(', ')':
			buf.UnreadRune()
			return s.String(), nil
		}

		s.WriteRune(r)
	}
}

func rangeParseQuotedValue(buf *bytes.Buffer) (string, error) {
	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case '\\':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
		case '"':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
			if r != '"' {
				buf.UnreadRune()
				return s.String(), nil
			}
		}
		s.WriteRune(r)
	}
}

type UntypedBinaryRange struct {
	Lower     []byte
	Upper     []byte
	LowerType BoundType
	UpperType BoundType
}

// 0 = ()      = 00000
// 1 = empty   = 00001
// 2 = [)      = 00010
// 4 = (]      = 00100
// 6 = []      = 00110
// 8 = )       = 01000
// 12 = ]      = 01100
// 16 = (      = 10000
// 18 = [      = 10010
// 24 =        = 11000

const emptyMask = 1
const lowerInclusiveMask = 2
const upperInclusiveMask = 4
const lowerUnboundedMask = 8
const upperUnboundedMask = 16

func ParseUntypedBinaryRange(src []byte) (*UntypedBinaryRange, error) {
	ubr := &UntypedBinaryRange{}

	if len(src) == 0 {
		return nil, errors.Errorf("range too short: %v", len(src))
	}

	rangeType := src[0]
	rp := 1

	if rangeType&emptyMask > 0 {
		if len(src[rp:]) > 0 {
			return nil, errors.Errorf("unexpected trailing bytes parsing empty range: %v", len(src[rp:]))
		}
		ubr.LowerType = Empty
		ubr.UpperType = Empty
		return ubr, nil
	}

	if rangeType&lowerInclusiveMask > 0 {
		ubr.LowerType = Inclusive
	} else if rangeType&lowerUnboundedMask > 0 {
		ubr.LowerType = Unbounded
	} else {
		ubr.LowerType = Exclusive
	}

	if rangeType&upperInclusiveMask > 0 {
		ubr.UpperType = Inclusive
	} else if rangeType&upperUnboundedMask > 0 {
		ubr.UpperType = Unbounded
	} else {
		ubr.UpperType = Exclusive
	}

	if ubr.LowerType == Unbounded && ubr.UpperType == Unbounded {
		if len(src[rp:]) > 0 {
			return nil, errors.Errorf("unexpected trailing bytes parsing unbounded range: %v", len(src[rp:]))
		}
		return ubr, nil
	}

	if len(src[rp:]) < 4 {
		return nil, errors.Errorf("too few bytes for size: %v", src[rp:])
	}
	valueLen := int(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	val := src[rp : rp+valueLen]
	rp += valueLen

	if ubr.LowerType != Unbounded {
		ubr.Lower = val
	} else {
		ubr.Upper = val
		if len(src[rp:]) > 0 {
			return nil, errors.Errorf("unexpected trailing bytes parsing range: %v", len(src[rp:]))
		}
		return ubr, nil
	}

	if ubr.UpperType != Unbounded {
		if len(src[rp:]) < 4 {
			return nil, errors.Errorf("too few bytes for size: %v", src[rp:])
		}
		valueLen := int(binary.BigEndian.Uint32(src[rp:]))
		rp += 4
		ubr.Upper = src[rp : rp+valueLen]
		rp += valueLen
	}

	if len(src[rp:]) > 0 {
		return nil, errors.Errorf("unexpected trailing bytes parsing range: %v", len(src[rp:]))
	}

	return ubr, nil

}
