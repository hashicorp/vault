// Copyright (c) 2018-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"fmt"
	"io"

	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	defaultChunkBufferSize  int64 = 8 << 10 // 8k
	defaultStringBufferSize int64 = 512
)

type largeChunkDecoder struct {
	r io.Reader

	rows  int // hint for number of rows
	cells int // hint for number of cells/row

	rem int // bytes remaining in rbuf
	ptr int // position in rbuf

	rbuf []byte
	sbuf *bytes.Buffer // buffer for decodeString

	ioError error
}

func decodeLargeChunk(r io.Reader, rowCount int, cellCount int) ([][]*string, error) {
	logger.Info("custom JSON Decoder")
	lcd := largeChunkDecoder{
		r, rowCount, cellCount,
		0, 0,
		make([]byte, defaultChunkBufferSize),
		bytes.NewBuffer(make([]byte, defaultStringBufferSize)),
		nil,
	}

	rows, err := lcd.decode()
	if lcd.ioError != nil && lcd.ioError != io.EOF {
		return nil, lcd.ioError
	} else if err != nil {
		return nil, err
	}

	return rows, nil
}

func (lcd *largeChunkDecoder) mkError(s string) error {
	return fmt.Errorf("corrupt chunk: %s", s)
}

func (lcd *largeChunkDecoder) decode() ([][]*string, error) {
	if lcd.nextByteNonWhitespace() != '[' {
		return nil, lcd.mkError("expected chunk to begin with '['")
	}

	rows := make([][]*string, 0, lcd.rows)
	if lcd.nextByteNonWhitespace() == ']' {
		return rows, nil // special case of an empty chunk
	}
	lcd.rewind(1)

OuterLoop:
	for {
		row, err := lcd.decodeRow()
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)

		switch c := lcd.nextByteNonWhitespace(); {
		case c == ',':
			continue // more elements in the array
		case c == ']':
			return rows, nil // we've scanned the whole chunk
		default:
			break OuterLoop
		}
	}
	return nil, lcd.mkError("invalid row boundary")
}

func (lcd *largeChunkDecoder) decodeRow() ([]*string, error) {
	if lcd.nextByteNonWhitespace() != '[' {
		return nil, lcd.mkError("expected row to begin with '['")
	}

	row := make([]*string, 0, lcd.cells)
	if lcd.nextByteNonWhitespace() == ']' {
		return row, nil // special case of an empty row
	}
	lcd.rewind(1)

OuterLoop:
	for {
		cell, err := lcd.decodeCell()
		if err != nil {
			return nil, err
		}
		row = append(row, cell)

		switch c := lcd.nextByteNonWhitespace(); {
		case c == ',':
			continue // more elements in the array
		case c == ']':
			return row, nil // we've scanned the whole row
		default:
			break OuterLoop
		}
	}
	return nil, lcd.mkError("invalid cell boundary")
}

func (lcd *largeChunkDecoder) decodeCell() (*string, error) {
	c := lcd.nextByteNonWhitespace()
	if c == '"' {
		s, err := lcd.decodeString()
		return &s, err
	} else if c == 'n' {
		if lcd.nextByte() == 'u' &&
			lcd.nextByte() == 'l' &&
			lcd.nextByte() == 'l' {
			return nil, nil
		}
	}
	return nil, lcd.mkError("cell begins with unexpected byte")
}

// TODO we can optimize this further by optimistically searching
// the read buffer for the next string. If it's short enough and
// doesn't contain any escaped characters, we can construct the
// return string directly without writing to the sbuf
func (lcd *largeChunkDecoder) decodeString() (string, error) {
	lcd.sbuf.Reset()
	for {
		// NOTE if you make changes here, ensure this
		// variable does not escape to the heap
		c := lcd.nextByte()
		if c == '"' {
			break
		} else if c == '\\' {
			if err := lcd.decodeEscaped(); err != nil {
				return "", err
			}
		} else if c < ' ' {
			return "", lcd.mkError("unexpected control character")
		} else if c < utf8.RuneSelf {
			lcd.sbuf.WriteByte(c)
		} else {
			lcd.rewind(1)
			lcd.sbuf.WriteRune(lcd.readRune())
		}
	}
	return lcd.sbuf.String(), nil
}

func (lcd *largeChunkDecoder) decodeEscaped() error {
	// NOTE if you make changes here, ensure this
	// variable does not escape to the heap
	c := lcd.nextByte()

	switch c {
	case '"', '\\', '/', '\'':
		lcd.sbuf.WriteByte(c)
	case 'b':
		lcd.sbuf.WriteByte('\b')
	case 'f':
		lcd.sbuf.WriteByte('\f')
	case 'n':
		lcd.sbuf.WriteByte('\n')
	case 'r':
		lcd.sbuf.WriteByte('\r')
	case 't':
		lcd.sbuf.WriteByte('\t')
	case 'u':
		rr := lcd.getu4()
		if rr < 0 {
			return lcd.mkError("invalid escape sequence")
		}
		if utf16.IsSurrogate(rr) {
			rr1, size := lcd.getu4WithPrefix()
			if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
				// A valid pair; consume.
				lcd.sbuf.WriteRune(dec)
				break
			}
			// Invalid surrogate; fall back to replacement rune.
			lcd.rewind(size)
			rr = unicode.ReplacementChar
		}
		lcd.sbuf.WriteRune(rr)
	default:
		return lcd.mkError("invalid escape sequence: " + string(c))
	}
	return nil
}

func (lcd *largeChunkDecoder) readRune() rune {
	lcd.ensureBytes(4)
	r, size := utf8.DecodeRune(lcd.rbuf[lcd.ptr:])
	lcd.ptr += size
	lcd.rem -= size
	return r
}

func (lcd *largeChunkDecoder) getu4WithPrefix() (rune, int) {
	lcd.ensureBytes(6)

	// NOTE take a snapshot of the cursor state. If this
	// is not a valid rune, then we need to roll back to
	// where we were before we began consuming bytes
	ptr := lcd.ptr

	if lcd.nextByte() != '\\' {
		return -1, lcd.ptr - ptr
	}
	if lcd.nextByte() != 'u' {
		return -1, lcd.ptr - ptr
	}
	r := lcd.getu4()
	return r, lcd.ptr - ptr
}

func (lcd *largeChunkDecoder) getu4() rune {
	var r rune
	for i := 0; i < 4; i++ {
		c := lcd.nextByte()
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

func (lcd *largeChunkDecoder) nextByteNonWhitespace() byte {
	for {
		c := lcd.nextByte()
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return c
		}
	}
}

func (lcd *largeChunkDecoder) rewind(n int) {
	lcd.ptr -= n
	lcd.rem += n
}

func (lcd *largeChunkDecoder) nextByte() byte {
	if lcd.rem == 0 {
		if lcd.ioError != nil {
			return 0
		}

		lcd.ptr = 0
		lcd.rem = lcd.fillBuffer(lcd.rbuf)
		if lcd.rem == 0 {
			return 0
		}
	}

	b := lcd.rbuf[lcd.ptr]
	lcd.ptr++

	lcd.rem--
	return b
}

func (lcd *largeChunkDecoder) ensureBytes(n int) {
	if lcd.rem <= n {
		rbuf := make([]byte, defaultChunkBufferSize)
		// NOTE when the buffer reads from the stream, there's no
		// guarantee that it will actually be filled. As such we
		// must use (ptr+rem) to compute the end of the slice.
		off := copy(rbuf, lcd.rbuf[lcd.ptr:lcd.ptr+lcd.rem])
		add := lcd.fillBuffer(rbuf[off:])

		lcd.ptr = 0
		lcd.rem += add
		lcd.rbuf = rbuf
	}
}

func (lcd *largeChunkDecoder) fillBuffer(b []byte) int {
	n, err := lcd.r.Read(b)
	if err != nil && err != io.EOF {
		lcd.ioError = err
		return 0
	} else if n <= 0 {
		lcd.ioError = io.EOF
		return 0
	}
	return n
}
