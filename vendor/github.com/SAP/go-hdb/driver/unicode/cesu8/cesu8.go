// Package cesu8 implements functions and constants to support text encoded in CESU-8.
// It implements functions comparable to the unicode/utf8 package for UTF-8 de- and encoding.
package cesu8

import (
	"unicode/utf16"
	"unicode/utf8"
)

const (
	// CESUMax is the maximum amount of bytes used by an CESU-8 codepoint encoding.
	CESUMax = 6
)

// Copied from unicode utf8.
const (
	tx = 0b10000000
	t3 = 0b11100000

	maskx = 0b00111111
	mask3 = 0b00001111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

// Size returns the amount of bytes needed to encode an UTF-8 byte slice to CESU-8.
func Size(p []byte) int {
	n := 0
	for len(p) > 0 {
		r, size := DecodeRune(p)
		n += RuneLen(r)
		p = p[size:]
	}
	return n
}

// StringSize is like Size with a string as parameter.
func StringSize(s string) int {
	n := 0
	for _, r := range s {
		n += RuneLen(r)
	}
	return n
}

// EncodeRune writes into p (which must be large enough) the CESU-8 encoding of the rune. It returns the number of bytes written.
func EncodeRune(p []byte, r rune) int {
	if r <= rune3Max {
		return utf8.EncodeRune(p, r)
	}
	high, low := utf16.EncodeRune(r)
	_ = p[5] // eliminate bounds checks
	p[0] = t3 | byte(high>>12)
	p[1] = tx | byte(high>>6)&maskx
	p[2] = tx | byte(high)&maskx
	p[3] = t3 | byte(low>>12)
	p[4] = tx | byte(low>>6)&maskx
	p[5] = tx | byte(low)&maskx
	return CESUMax
}

// FullRune reports whether the bytes in p begin with a full CESU-8 encoding of a rune.
func FullRune(p []byte) bool {
	if isSurrogate(p) {
		return isSurrogate(p[3:])
	}
	return utf8.FullRune(p)
}

// DecodeRune unpacks the first CESU-8 encoding in p and returns the rune and its width in bytes.
func DecodeRune(p []byte) (rune, int) {
	if !isSurrogate(p) {
		return utf8.DecodeRune(p)
	}
	high := decodeCheckedSurrogate(p)
	low, ok := decodeSurrogate(p[3:])
	if !ok {
		return utf8.RuneError, 3
	}
	return utf16.DecodeRune(high, low), CESUMax
}

// RuneLen returns the number of bytes required to encode the rune.
func RuneLen(r rune) int {
	switch {
	case r < 0:
		return -1
	case r <= rune1Max:
		return 1
	case r <= rune2Max:
		return 2
	case r <= rune3Max:
		return 3
	case r <= utf8.MaxRune:
		return CESUMax
	default:
		return -1
	}
}

const (
	sp0    = 0xed
	sb1Min = 0xa0
	sb1Max = 0xbf
)

func decodeSurrogate(p []byte) (rune, bool) {
	if len(p) < 3 {
		return utf8.RuneError, false
	}
	p0 := p[0]
	if p0 != sp0 {
		return utf8.RuneError, false
	}
	b1 := p[1]
	if b1 < sb1Min || b1 > sb1Max {
		return utf8.RuneError, false
	}
	b2 := p[2]
	return rune(p0&mask3)<<12 | rune(b1&maskx)<<6 | rune(b2&maskx), true
}

func decodeCheckedSurrogate(p []byte) rune {
	return rune(p[0]&mask3)<<12 | rune(p[1]&maskx)<<6 | rune(p[2]&maskx)
}

func isSurrogate(p []byte) bool {
	if len(p) < 3 {
		return false
	}
	b1 := p[1]
	if p[0] != sp0 || b1 < sb1Min || b1 > sb1Max {
		return false
	}
	return true
}
