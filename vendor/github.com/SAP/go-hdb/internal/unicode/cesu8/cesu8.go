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

// Size returns the amount of bytes needed to encode an UTF-8 byte slice to CESU-8.
func Size(p []byte) int {
	n := 0
	for i := 0; i < len(p); {
		r, size, _ := decodeRune(p[i:])
		i += size
		n += RuneLen(r)
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
		return encodeRune(p, r)
	}
	high, low := utf16.EncodeRune(r)
	n := encodeRune(p, high)
	n += encodeRune(p[n:], low)
	return n
}

// FullRune reports whether the bytes in p begin with a full CESU-8 encoding of a rune.
func FullRune(p []byte) bool {
	high, n, short := decodeRune(p)
	if short {
		return false
	}
	if !utf16.IsSurrogate(high) {
		return true
	}
	_, _, short = decodeRune(p[n:])
	return !short
}

// DecodeRune unpacks the first CESU-8 encoding in p and returns the rune and its width in bytes.
func DecodeRune(p []byte) (rune, int) {
	high, n1, _ := decodeRune(p)
	if !utf16.IsSurrogate(high) {
		return high, n1
	}
	low, n2, _ := decodeRune(p[n1:])
	if low == utf8.RuneError {
		return low, n1 + n2
	}
	return utf16.DecodeRune(high, low), n1 + n2
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
	}
	return -1
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// Copied from unicode utf8
// - allow utf8 encoding of utf16 surrogate values
// - see (*) for code changes

// Code points in the surrogate range are not valid for UTF-8.
const (
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1
)

func encodeRune(p []byte, r rune) int {
	// Negative values are erroneous.  Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		p[0] = byte(r)
		return 1
	case i <= rune2Max:
		p[0] = t2 | byte(r>>6)
		p[1] = tx | byte(r)&maskx
		return 2
	//case i > MaxRune, surrogateMin <= i && i <= surrogateMax: // replaced (*)
	case i > utf8.MaxRune: // (*)
		r = utf8.RuneError
		fallthrough
	case i <= rune3Max:
		p[0] = t3 | byte(r>>12)
		p[1] = tx | byte(r>>6)&maskx
		p[2] = tx | byte(r)&maskx
		return 3
	default:
		p[0] = t4 | byte(r>>18)
		p[1] = tx | byte(r>>12)&maskx
		p[2] = tx | byte(r>>6)&maskx
		p[3] = tx | byte(r)&maskx
		return 4
	}
}

func decodeRune(p []byte) (r rune, size int, short bool) {
	n := len(p)
	if n < 1 {
		return utf8.RuneError, 0, true
	}
	c0 := p[0]

	// 1-byte, 7-bit sequence?
	if c0 < tx {
		return rune(c0), 1, false
	}

	// unexpected continuation byte?
	if c0 < t2 {
		return utf8.RuneError, 1, false
	}

	// need first continuation byte
	if n < 2 {
		return utf8.RuneError, 1, true
	}
	c1 := p[1]
	if c1 < tx || t2 <= c1 {
		return utf8.RuneError, 1, false
	}

	// 2-byte, 11-bit sequence?
	if c0 < t3 {
		r = rune(c0&mask2)<<6 | rune(c1&maskx)
		if r <= rune1Max {
			return utf8.RuneError, 1, false
		}
		return r, 2, false
	}

	// need second continuation byte
	if n < 3 {
		return utf8.RuneError, 1, true
	}
	c2 := p[2]
	if c2 < tx || t2 <= c2 {
		return utf8.RuneError, 1, false
	}

	// 3-byte, 16-bit sequence?
	if c0 < t4 {
		r = rune(c0&mask3)<<12 | rune(c1&maskx)<<6 | rune(c2&maskx)
		if r <= rune2Max {
			return utf8.RuneError, 1, false
		}
		// do not throw error on surrogates // (*)
		//if surrogateMin <= r && r <= surrogateMax {
		//	return RuneError, 1, false
		//}
		return r, 3, false
	}

	// need third continuation byte
	if n < 4 {
		return utf8.RuneError, 1, true
	}
	c3 := p[3]
	if c3 < tx || t2 <= c3 {
		return utf8.RuneError, 1, false
	}

	// 4-byte, 21-bit sequence?
	if c0 < t5 {
		r = rune(c0&mask4)<<18 | rune(c1&maskx)<<12 | rune(c2&maskx)<<6 | rune(c3&maskx)
		if r <= rune3Max || utf8.MaxRune < r {
			return utf8.RuneError, 1, false
		}
		return r, 4, false
	}

	// error
	return utf8.RuneError, 1, false
}
