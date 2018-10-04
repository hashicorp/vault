package pgtype

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

// Information on the internals of PostgreSQL arrays can be found in
// src/include/utils/array.h and src/backend/utils/adt/arrayfuncs.c. Of
// particular interest is the array_send function.

type ArrayHeader struct {
	ContainsNull bool
	ElementOID   int32
	Dimensions   []ArrayDimension
}

type ArrayDimension struct {
	Length     int32
	LowerBound int32
}

func (dst *ArrayHeader) DecodeBinary(ci *ConnInfo, src []byte) (int, error) {
	if len(src) < 12 {
		return 0, errors.Errorf("array header too short: %d", len(src))
	}

	rp := 0

	numDims := int(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	dst.ContainsNull = binary.BigEndian.Uint32(src[rp:]) == 1
	rp += 4

	dst.ElementOID = int32(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	if numDims > 0 {
		dst.Dimensions = make([]ArrayDimension, numDims)
	}
	if len(src) < 12+numDims*8 {
		return 0, errors.Errorf("array header too short for %d dimensions: %d", numDims, len(src))
	}
	for i := range dst.Dimensions {
		dst.Dimensions[i].Length = int32(binary.BigEndian.Uint32(src[rp:]))
		rp += 4

		dst.Dimensions[i].LowerBound = int32(binary.BigEndian.Uint32(src[rp:]))
		rp += 4
	}

	return rp, nil
}

func (src *ArrayHeader) EncodeBinary(ci *ConnInfo, buf []byte) []byte {
	buf = pgio.AppendInt32(buf, int32(len(src.Dimensions)))

	var containsNull int32
	if src.ContainsNull {
		containsNull = 1
	}
	buf = pgio.AppendInt32(buf, containsNull)

	buf = pgio.AppendInt32(buf, src.ElementOID)

	for i := range src.Dimensions {
		buf = pgio.AppendInt32(buf, src.Dimensions[i].Length)
		buf = pgio.AppendInt32(buf, src.Dimensions[i].LowerBound)
	}

	return buf
}

type UntypedTextArray struct {
	Elements   []string
	Dimensions []ArrayDimension
}

func ParseUntypedTextArray(src string) (*UntypedTextArray, error) {
	dst := &UntypedTextArray{}

	buf := bytes.NewBufferString(src)

	skipWhitespace(buf)

	r, _, err := buf.ReadRune()
	if err != nil {
		return nil, errors.Errorf("invalid array: %v", err)
	}

	var explicitDimensions []ArrayDimension

	// Array has explicit dimensions
	if r == '[' {
		buf.UnreadRune()

		for {
			r, _, err = buf.ReadRune()
			if err != nil {
				return nil, errors.Errorf("invalid array: %v", err)
			}

			if r == '=' {
				break
			} else if r != '[' {
				return nil, errors.Errorf("invalid array, expected '[' or '=' got %v", r)
			}

			lower, err := arrayParseInteger(buf)
			if err != nil {
				return nil, errors.Errorf("invalid array: %v", err)
			}

			r, _, err = buf.ReadRune()
			if err != nil {
				return nil, errors.Errorf("invalid array: %v", err)
			}

			if r != ':' {
				return nil, errors.Errorf("invalid array, expected ':' got %v", r)
			}

			upper, err := arrayParseInteger(buf)
			if err != nil {
				return nil, errors.Errorf("invalid array: %v", err)
			}

			r, _, err = buf.ReadRune()
			if err != nil {
				return nil, errors.Errorf("invalid array: %v", err)
			}

			if r != ']' {
				return nil, errors.Errorf("invalid array, expected ']' got %v", r)
			}

			explicitDimensions = append(explicitDimensions, ArrayDimension{LowerBound: lower, Length: upper - lower + 1})
		}

		r, _, err = buf.ReadRune()
		if err != nil {
			return nil, errors.Errorf("invalid array: %v", err)
		}
	}

	if r != '{' {
		return nil, errors.Errorf("invalid array, expected '{': %v", err)
	}

	implicitDimensions := []ArrayDimension{{LowerBound: 1, Length: 0}}

	// Consume all initial opening brackets. This provides number of dimensions.
	for {
		r, _, err = buf.ReadRune()
		if err != nil {
			return nil, errors.Errorf("invalid array: %v", err)
		}

		if r == '{' {
			implicitDimensions[len(implicitDimensions)-1].Length = 1
			implicitDimensions = append(implicitDimensions, ArrayDimension{LowerBound: 1})
		} else {
			buf.UnreadRune()
			break
		}
	}
	currentDim := len(implicitDimensions) - 1
	counterDim := currentDim

	for {
		r, _, err = buf.ReadRune()
		if err != nil {
			return nil, errors.Errorf("invalid array: %v", err)
		}

		switch r {
		case '{':
			if currentDim == counterDim {
				implicitDimensions[currentDim].Length++
			}
			currentDim++
		case ',':
		case '}':
			currentDim--
			if currentDim < counterDim {
				counterDim = currentDim
			}
		default:
			buf.UnreadRune()
			value, err := arrayParseValue(buf)
			if err != nil {
				return nil, errors.Errorf("invalid array value: %v", err)
			}
			if currentDim == counterDim {
				implicitDimensions[currentDim].Length++
			}
			dst.Elements = append(dst.Elements, value)
		}

		if currentDim < 0 {
			break
		}
	}

	skipWhitespace(buf)

	if buf.Len() > 0 {
		return nil, errors.Errorf("unexpected trailing data: %v", buf.String())
	}

	if len(dst.Elements) == 0 {
		dst.Dimensions = nil
	} else if len(explicitDimensions) > 0 {
		dst.Dimensions = explicitDimensions
	} else {
		dst.Dimensions = implicitDimensions
	}

	return dst, nil
}

func skipWhitespace(buf *bytes.Buffer) {
	var r rune
	var err error
	for r, _, _ = buf.ReadRune(); unicode.IsSpace(r); r, _, _ = buf.ReadRune() {
	}

	if err != io.EOF {
		buf.UnreadRune()
	}
}

func arrayParseValue(buf *bytes.Buffer) (string, error) {
	r, _, err := buf.ReadRune()
	if err != nil {
		return "", err
	}
	if r == '"' {
		return arrayParseQuotedValue(buf)
	}
	buf.UnreadRune()

	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case ',', '}':
			buf.UnreadRune()
			return s.String(), nil
		}

		s.WriteRune(r)
	}
}

func arrayParseQuotedValue(buf *bytes.Buffer) (string, error) {
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
			buf.UnreadRune()
			return s.String(), nil
		}
		s.WriteRune(r)
	}
}

func arrayParseInteger(buf *bytes.Buffer) (int32, error) {
	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return 0, err
		}

		if '0' <= r && r <= '9' {
			s.WriteRune(r)
		} else {
			buf.UnreadRune()
			n, err := strconv.ParseInt(s.String(), 10, 32)
			if err != nil {
				return 0, err
			}
			return int32(n), nil
		}
	}
}

func EncodeTextArrayDimensions(buf []byte, dimensions []ArrayDimension) []byte {
	var customDimensions bool
	for _, dim := range dimensions {
		if dim.LowerBound != 1 {
			customDimensions = true
		}
	}

	if !customDimensions {
		return buf
	}

	for _, dim := range dimensions {
		buf = append(buf, '[')
		buf = append(buf, strconv.FormatInt(int64(dim.LowerBound), 10)...)
		buf = append(buf, ':')
		buf = append(buf, strconv.FormatInt(int64(dim.LowerBound+dim.Length-1), 10)...)
		buf = append(buf, ']')
	}

	return append(buf, '=')
}

var quoteArrayReplacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func quoteArrayElement(src string) string {
	return `"` + quoteArrayReplacer.Replace(src) + `"`
}

func QuoteArrayElementIfNeeded(src string) string {
	if src == "" || (len(src) == 4 && strings.ToLower(src) == "null") || src[0] == ' ' || src[len(src)-1] == ' ' || strings.ContainsAny(src, `{},"\`) {
		return quoteArrayElement(src)
	}
	return src
}
