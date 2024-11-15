// Package cesu8 implements functions and constants to support text encoded in CESU-8.
// It implements functions comparable to the unicode/utf8 package for UTF-8 de- and encoding.
package cesu8

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// Encoding constants.
const (
	UTF8  = "UTF-8"
	CESU8 = "CESU-8"
)

// DecodeError is raised when a transformer detects invalid encoded data.
type DecodeError struct {
	enc string // encoding
	p   int    // position of error in value
	v   []byte // value
}

func newDecodeError(enc string, p int, v []byte) *DecodeError {
	// copy value
	cv := make([]byte, len(v))
	copy(cv, v)
	return &DecodeError{enc: enc, p: p, v: cv}
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("invalid %s: %x at position %d", e.enc, e.v, e.p)
}

// Enc returns the expected encoding of the erroneous data.
func (e *DecodeError) Enc() string { return e.enc }

// Pos returns the position of the invalid rune.
func (e *DecodeError) Pos() int { return e.p }

// Value returns the value which should be decoded.
func (e *DecodeError) Value() []byte { return e.v }

// Encoder supports encoding of UTF-8 encoded data into CESU-8.
type Encoder struct {
	transform.NopResetter
	errorHandler func(err *DecodeError) (rune, error)
}

// NewEncoder creates a new encoder instance. With parameter errorHandler a custom error handling function could be used in case
// the encoder would detect invalid UTF-8 encoded characters.
func NewEncoder(errorHandler func(err *DecodeError) (rune, error)) *Encoder {
	return &Encoder{errorHandler: errorHandler}
}

// Transform implements the transform.Transformer interface.
func (e *Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	i, j := 0, 0
	for i < len(src) {
		if src[i] < utf8.RuneSelf {
			if j >= len(dst) {
				return j, i, transform.ErrShortDst
			}
			dst[j] = src[i]
			i++
			j++
			continue
		}
		// check if additional bytes needed (ErrShortSrc) only
		// - if further bytes are potentially available (!atEOF) and
		// - remaining buffer smaller than max size for an encoded UTF-8 rune
		if !atEOF && len(src[i:]) < utf8.UTFMax {
			if !utf8.FullRune(src[i:]) {
				return j, i, transform.ErrShortSrc
			}
		}
		r, n := utf8.DecodeRune(src[i:])
		if r == utf8.RuneError {
			decodeErr := newDecodeError(UTF8, i, src)
			if e.errorHandler == nil {
				return j, i, decodeErr
			}
			r, err = e.errorHandler(decodeErr)
			if err != nil {
				return j, i, err
			}
		}
		m := RuneLen(r)
		switch {
		case m == -1:
			panic("internal UTF-8 to CESU-8 transformation error")
		case j+m > len(dst):
			return j, i, transform.ErrShortDst
		}
		EncodeRune(dst[j:], r)
		i += n
		j += m
	}
	return j, i, nil
}

// Decoder supports decoding of CESU-8 encoded data into UTF-8.
type Decoder struct {
	transform.NopResetter
	errorHandler func(err *DecodeError) (rune, error)
}

// NewDecoder creates a new decoder instance. With parameter errorHandler a custom error handling function could be used in case
// the decoder would detect invalid CESU-8 encoded characters.
func NewDecoder(errorHandler func(err *DecodeError) (rune, error)) *Decoder {
	return &Decoder{errorHandler: errorHandler}
}

// Transform implements the transform.Transformer interface.
func (d *Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	i, j := 0, 0
	for i < len(src) {
		if src[i] < utf8.RuneSelf {
			if j >= len(dst) {
				return j, i, transform.ErrShortDst
			}
			dst[j] = src[i]
			i++
			j++
			continue
		}
		// check if additional bytes needed (ErrShortSrc) only
		// - if further bytes are potentially available (!atEOF) and
		// - remaining buffer smaller than max size for an encoded CESU-8 rune
		if !atEOF && len(src[i:]) < CESUMax {
			if !FullRune(src[i:]) {
				return j, i, transform.ErrShortSrc
			}
		}
		r, n := DecodeRune(src[i:])
		if r == utf8.RuneError {
			decodeErr := newDecodeError(CESU8, i, src)
			if d.errorHandler == nil {
				return j, i, decodeErr
			}
			r, err = d.errorHandler(decodeErr)
			if err != nil {
				return j, i, err
			}
		}
		m := utf8.RuneLen(r)
		switch {
		case m == -1:
			panic("internal CESU-8 to UTF-8 transformation error")
		case j+m > len(dst):
			return j, i, transform.ErrShortDst
		}
		utf8.EncodeRune(dst[j:], r)
		i += n
		j += m
	}
	return j, i, nil
}

var (
	defaultDecoder = NewDecoder(nil)
	defaultEncoder = NewEncoder(nil)
)

// DefaultDecoder returns the default CESU-8 to UTF-8 decoder.
func DefaultDecoder() transform.Transformer { return defaultDecoder }

// DefaultEncoder returns the default UTF-8 to CESU-8 encoder.
func DefaultEncoder() transform.Transformer { return defaultEncoder }

// ReplaceErrorHandler is a decoding error handling function replacing invalid CESU-8 data with the
// unicode replacement character '\uFFFD'.
func ReplaceErrorHandler(err *DecodeError) (rune, error) { return unicode.ReplacementChar, nil }
