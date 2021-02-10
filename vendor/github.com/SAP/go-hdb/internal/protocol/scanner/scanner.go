// +build go1.10

// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

/*
Package scanner implements a HANA SQL query scanner.

For a detailed HANA SQL query syntax please see
https://help.sap.com/doc/6254b3bb439c4f409a979dc407b49c9b/2.0.00/en-US/SAP_HANA_SQL_Script_Reference_en.pdf
*/
package scanner

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ErrToken is raised when a token is malformed (e.g. string with missing ending quote).
var ErrToken = errors.New("invalid token")

// Token is the set of lexicals tokens of a SQL query.
type Token int

// Token constants.
const (
	EOS Token = iota
	Error
	Undefined
	Operator
	Delimiter
	IdentifierDelimiter
	Identifier
	QuotedIdentifier
	Variable
	PosVariable
	NamedVariable
	String
	Number
)

var tokenString = map[Token]string{
	EOS:                 "EndOfString",
	Error:               "Error",
	Undefined:           "Undefined",
	Operator:            "Operator",
	Delimiter:           "Delimiter",
	IdentifierDelimiter: "IdentifierDelimiter",
	Identifier:          "Identifier",
	QuotedIdentifier:    "QuotedIdentifier",
	Variable:            "Variable",
	PosVariable:         "PosVariable",
	NamedVariable:       "NamedVariable",
	String:              "String",
	Number:              "Number",
}

func (t Token) String() string {
	if s, ok := tokenString[t]; ok {
		return s
	}
	return fmt.Sprintf("%d", t)
}

var compositeOperators = map[string]struct{}{"<>": {}, "<=": {}, ">=": {}, "!=": {}}

func isOperator(ch rune) bool           { return strings.ContainsRune("<>=!", ch) }
func isCompositeOperator(s string) bool { _, ok := compositeOperators[s]; return ok }
func isDelimiter(ch rune) bool          { return strings.ContainsRune(",;(){}[]", ch) }
func isNameDelimiter(ch rune) bool      { return ch == '.' }
func isDigit(ch rune) bool              { return unicode.IsDigit(ch) }
func isNumber(ch rune) bool             { return ch == '+' || ch == '-' || isDigit(ch) }
func isExp(ch rune) bool                { return ch == 'e' || ch == 'E' }
func isDecimalSeparator(ch rune) bool   { return ch == '.' }
func isIdentifier(ch rune) bool         { return ch == '_' || unicode.IsLetter(ch) }
func isAlpha(ch rune) bool              { return ch == '#' || ch == '$' || isIdentifier(ch) || isDigit(ch) }
func isSingleQuote(ch rune) bool        { return ch == '\'' }
func isDoubleQuote(ch rune) bool        { return ch == '"' }
func isQuestionMark(ch rune) bool       { return ch == '?' }
func isColon(ch rune) bool              { return ch == ':' }

// A Scanner implements reading of SQL query tokens.
type Scanner struct {
	s    string
	i    int // reading position
	prev int // previous reading position
}

// Reset initializes a Scanner with a new SQL statement.
func (sc *Scanner) Reset(s string) {
	sc.s = s
	sc.i = 0
	sc.prev = -1
}

func (sc *Scanner) readRune() (rune, bool) {
	sc.prev = sc.i
	if sc.i >= len(sc.s) {
		return 0, false
	}
	if c := sc.s[sc.i]; c < utf8.RuneSelf {
		sc.i++
		return rune(c), true
	}
	ch, size := utf8.DecodeRuneInString(sc.s[sc.i:])
	sc.i += size
	return ch, true
}

func (sc *Scanner) unreadRune() {
	if sc.prev == -1 {
		panic("unreadRune before readRune")
	}
	sc.i = sc.prev
	sc.prev = -1
}

func (sc *Scanner) scanWhitespace() {
	for {
		ch, ok := sc.readRune()
		if !ok {
			return
		}
		if !unicode.IsSpace(ch) {
			sc.unreadRune()
			return
		}
	}
}

func (sc *Scanner) scanOperator(ch rune) {
	ch2, ok := sc.readRune()
	if !ok {
		return
	}
	if isCompositeOperator(string([]rune{ch, ch2})) {
		return
	}
	sc.unreadRune()
}

func (sc *Scanner) scanNumeric() {
	for {
		ch, ok := sc.readRune()
		if !ok {
			return
		}
		if !isDigit(ch) {
			sc.unreadRune()
			return
		}
	}
}

func (sc *Scanner) scanAlpha() {
	for {
		ch, ok := sc.readRune()
		if !ok {
			return
		}
		if !isAlpha(ch) {
			sc.unreadRune()
			return
		}
	}
}

func (sc *Scanner) scanQuotedIdentifier(quote rune) Token {
	for {
		ch, ok := sc.readRune()
		if !ok {
			return Error
		}
		if ch == quote {
			ch, ok := sc.readRune()
			if !ok {
				return QuotedIdentifier
			}
			if ch != quote {
				sc.unreadRune()
				return QuotedIdentifier
			}
		}
	}
}

func (sc *Scanner) scanVariable() Token {
	ch, ok := sc.readRune()
	if !ok {
		return Error
	}
	if isDigit(ch) {
		sc.scanNumeric()
		return PosVariable
	}
	sc.scanAlpha()
	return NamedVariable
}

func (sc *Scanner) scanNumber() Token {
	sc.scanNumeric()
	ch, ok := sc.readRune()
	if !ok {
		return Number
	}
	if isDecimalSeparator(ch) {
		sc.scanNumeric()
	}
	ch, ok = sc.readRune()
	if !ok {
		return Number
	}
	if isExp(ch) {
		if isNumber(ch) {
			sc.scanNumeric()
			return Number
		}
		return Error
	}
	return Number
}

// Next reads and returns the next token.
func (sc *Scanner) Next() (token Token, start, end int) {
	sc.scanWhitespace()

	start = sc.i

	ch, ok := sc.readRune()
	if !ok {
		return EOS, start, sc.i
	}

	switch {
	default:
		return Error, start, sc.i

	case isDelimiter(ch):
		return Delimiter, start, sc.i
	case isNameDelimiter(ch):
		return IdentifierDelimiter, start, sc.i
	case isQuestionMark(ch):
		return Variable, start, sc.i

	case isOperator(ch):
		sc.scanOperator(ch)
		return Operator, start, sc.i

	case isIdentifier(ch):
		sc.scanAlpha()
		return Identifier, start, sc.i

	case isSingleQuote(ch) || isDoubleQuote(ch):
		return sc.scanQuotedIdentifier(ch), start, sc.i

	case isColon(ch):
		return sc.scanVariable(), start, sc.i

	case isNumber(ch):
		return sc.scanNumber(), start, sc.i
	}
}
