package ini

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	runesTrue  = []rune("true")
	runesFalse = []rune("false")
)

var literalValues = [][]rune{
	runesTrue,
	runesFalse,
}

func isBoolValue(b []rune) bool {
	for _, lv := range literalValues {
		if isCaselessLitValue(lv, b) {
			return true
		}
	}
	return false
}

func isLitValue(want, have []rune) bool {
	if len(have) < len(want) {
		return false
	}

	for i := 0; i < len(want); i++ {
		if want[i] != have[i] {
			return false
		}
	}

	return true
}

// isCaselessLitValue is a caseless value comparison, assumes want is already lower-cased for efficiency.
func isCaselessLitValue(want, have []rune) bool {
	if len(have) < len(want) {
		return false
	}

	for i := 0; i < len(want); i++ {
		if want[i] != unicode.ToLower(have[i]) {
			return false
		}
	}

	return true
}

// isNumberValue will return whether not the leading characters in
// a byte slice is a number. A number is delimited by whitespace or
// the newline token.
//
// A number is defined to be in a binary, octal, decimal (int | float), hex format,
// or in scientific notation.
func isNumberValue(b []rune) bool {
	negativeIndex := 0
	helper := numberHelper{}
	needDigit := false

	for i := 0; i < len(b); i++ {
		negativeIndex++

		switch b[i] {
		case '-':
			if helper.IsNegative() || negativeIndex != 1 {
				return false
			}
			helper.Determine(b[i])
			needDigit = true
			continue
		case 'e', 'E':
			if err := helper.Determine(b[i]); err != nil {
				return false
			}
			negativeIndex = 0
			needDigit = true
			continue
		case 'b':
			if helper.numberFormat == hex {
				break
			}
			fallthrough
		case 'o', 'x':
			needDigit = true
			if i == 0 {
				return false
			}

			fallthrough
		case '.':
			if err := helper.Determine(b[i]); err != nil {
				return false
			}
			needDigit = true
			continue
		}

		if i > 0 && (isNewline(b[i:]) || isWhitespace(b[i])) {
			return !needDigit
		}

		if !helper.CorrectByte(b[i]) {
			return false
		}
		needDigit = false
	}

	return !needDigit
}

func isValid(b []rune) (bool, int, error) {
	if len(b) == 0 {
		// TODO: should probably return an error
		return false, 0, nil
	}

	return isValidRune(b[0]), 1, nil
}

func isValidRune(r rune) bool {
	return r != ':' && r != '=' && r != '[' && r != ']' && r != ' ' && r != '\n'
}

// ValueType is an enum that will signify what type
// the Value is
type ValueType int

func (v ValueType) String() string {
	switch v {
	case NoneType:
		return "NONE"
	case DecimalType:
		return "FLOAT"
	case IntegerType:
		return "INT"
	case StringType:
		return "STRING"
	case BoolType:
		return "BOOL"
	}

	return ""
}

// ValueType enums
const (
	NoneType = ValueType(iota)
	DecimalType // deprecated
	IntegerType // deprecated
	StringType
	QuotedStringType
	BoolType // deprecated
)

// Value is a union container
type Value struct {
	Type ValueType
	raw  []rune

	integer int64 // deprecated
	decimal float64 // deprecated
	boolean bool // deprecated
	str     string
}

func newValue(t ValueType, base int, raw []rune) (Value, error) {
	v := Value{
		Type: t,
		raw:  raw,
	}
	var err error

	switch t {
	case DecimalType:
		v.decimal, err = strconv.ParseFloat(string(raw), 64)
	case IntegerType:
		if base != 10 {
			raw = raw[2:]
		}

		v.integer, err = strconv.ParseInt(string(raw), base, 64)
	case StringType:
		v.str = string(raw)
	case QuotedStringType:
		v.str = string(raw[1 : len(raw)-1])
	case BoolType:
		v.boolean = isCaselessLitValue(runesTrue, v.raw)
	}

	// issue 2253
	//
	// if the value trying to be parsed is too large, then we will use
	// the 'StringType' and raw value instead.
	if nerr, ok := err.(*strconv.NumError); ok && nerr.Err == strconv.ErrRange {
		v.Type = StringType
		v.str = string(raw)
		err = nil
	}

	return v, err
}

// Append will append values and change the type to a string
// type.
func (v *Value) Append(tok Token) {
	r := tok.Raw()
	if v.Type != QuotedStringType {
		v.Type = StringType
		r = tok.raw[1 : len(tok.raw)-1]
	}
	if tok.Type() != TokenLit {
		v.raw = append(v.raw, tok.Raw()...)
	} else {
		v.raw = append(v.raw, r...)
	}
}

func (v Value) String() string {
	switch v.Type {
	case DecimalType:
		return fmt.Sprintf("decimal: %f", v.decimal)
	case IntegerType:
		return fmt.Sprintf("integer: %d", v.integer)
	case StringType:
		return fmt.Sprintf("string: %s", string(v.raw))
	case QuotedStringType:
		return fmt.Sprintf("quoted string: %s", string(v.raw))
	case BoolType:
		return fmt.Sprintf("bool: %t", v.boolean)
	default:
		return "union not set"
	}
}

func newLitToken(b []rune) (Token, int, error) {
	n := 0
	var err error

	token := Token{}
	if b[0] == '"' {
		n, err = getStringValue(b)
		if err != nil {
			return token, n, err
		}

		token = newToken(TokenLit, b[:n], QuotedStringType)
	} else {
		n, err = getValue(b)
		token = newToken(TokenLit, b[:n], StringType)
	}

	return token, n, err
}

// IntValue returns an integer value
func (v Value) IntValue() (int64, bool) {
	i, err := strconv.ParseInt(string(v.raw), 0, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

// FloatValue returns a float value
func (v Value) FloatValue() (float64, bool) {
	f, err := strconv.ParseFloat(string(v.raw), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

// BoolValue returns a bool value
func (v Value) BoolValue() (bool, bool) {
	// we don't use ParseBool as it recognizes more than what we've
	// historically supported
	if isCaselessLitValue(runesTrue, v.raw) {
		return true, true
	} else if isCaselessLitValue(runesFalse, v.raw) {
		return false, true
	}
	return false, false
}

func isTrimmable(r rune) bool {
	switch r {
	case '\n', ' ':
		return true
	}
	return false
}

// StringValue returns the string value
func (v Value) StringValue() string {
	switch v.Type {
	case StringType:
		return strings.TrimFunc(string(v.raw), isTrimmable)
	case QuotedStringType:
		// preserve all characters in the quotes
		return string(removeEscapedCharacters(v.raw[1 : len(v.raw)-1]))
	default:
		return strings.TrimFunc(string(v.raw), isTrimmable)
	}
}

func contains(runes []rune, c rune) bool {
	for i := 0; i < len(runes); i++ {
		if runes[i] == c {
			return true
		}
	}

	return false
}

func runeCompare(v1 []rune, v2 []rune) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}
