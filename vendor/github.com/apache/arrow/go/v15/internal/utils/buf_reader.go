// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// bufferedReader is similar to bufio.Reader except
// it will expand the buffer if necessary when asked to Peek
// more bytes than are in the buffer
type bufferedReader struct {
	bufferSz int
	buf      []byte
	r, w     int
	rd       io.Reader
	err      error
}

// NewBufferedReader returns a buffered reader with similar semantics to bufio.Reader
// except Peek will expand the internal buffer if needed rather than return
// an error.
func NewBufferedReader(rd io.Reader, sz int) *bufferedReader {
	// if rd is already a buffered reader whose buffer is >= the requested size
	// then just return it as is. no need to make a new object.
	b, ok := rd.(*bufferedReader)
	if ok && len(b.buf) >= sz {
		return b
	}

	r := &bufferedReader{
		rd: rd,
	}
	r.resizeBuffer(sz)
	return r
}

func (b *bufferedReader) resetBuffer() {
	if b.buf == nil {
		b.buf = make([]byte, b.bufferSz)
	} else if b.bufferSz > cap(b.buf) {
		buf := b.buf
		b.buf = make([]byte, b.bufferSz)
		copy(b.buf, buf)
	} else {
		b.buf = b.buf[:b.bufferSz]
	}
}

func (b *bufferedReader) resizeBuffer(newSize int) {
	b.bufferSz = newSize
	b.resetBuffer()
}

func (b *bufferedReader) fill() error {
	// slide existing data to the beginning
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		return fmt.Errorf("arrow/bufferedreader: %w", bufio.ErrBufferFull)
	}

	n, err := io.ReadAtLeast(b.rd, b.buf[b.w:], 1)
	if n < 0 {
		return fmt.Errorf("arrow/bufferedreader: filling buffer: %w", bufio.ErrNegativeCount)
	}

	b.w += n
	b.err = err
	return nil
}

func (b *bufferedReader) readErr() error {
	err := b.err
	b.err = nil
	return err
}

// Buffered returns the number of bytes currently buffered
func (b *bufferedReader) Buffered() int { return b.w - b.r }

// SetBufferSize resets the size of the internal buffer to the desired size.
// Will return an error if newSize is <= 0 or if newSize is less than the size
// of the buffered data.
func (b *bufferedReader) SetBufferSize(newSize int) error {
	if newSize <= 0 {
		return errors.New("buffer size should be positive")
	}

	if b.w >= newSize {
		return errors.New("cannot shrink read buffer if buffered data remains")
	}

	b.resizeBuffer(newSize)
	return nil
}

// Peek will buffer and return n bytes from the underlying reader without advancing
// the reader itself. If n is larger than the current buffer size, the buffer will
// be expanded to accommodate the extra bytes rather than error.
func (b *bufferedReader) Peek(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("arrow/bufferedreader: %w", bufio.ErrNegativeCount)
	}

	if n > len(b.buf) {
		if err := b.SetBufferSize(n); err != nil {
			return nil, err
		}
	}

	for b.w-b.r < n && b.w-b.r < len(b.buf) && b.err == nil {
		b.fill() // b.w-b.r < len(b.buf) => buffer is not full
	}

	return b.buf[b.r : b.r+n], b.readErr()
}

// Discard skips the next n bytes either by advancing the internal buffer
// or by reading that many bytes in and throwing them away.
func (b *bufferedReader) Discard(n int) (discarded int, err error) {
	if n < 0 {
		return 0, fmt.Errorf("arrow/bufferedreader: %w", bufio.ErrNegativeCount)
	}

	if n == 0 {
		return
	}

	remain := n
	for {
		skip := b.Buffered()
		if skip == 0 {
			b.fill()
			skip = b.Buffered()
		}
		if skip > remain {
			skip = remain
		}
		b.r += skip
		remain -= skip
		if remain == 0 {
			return n, nil
		}
		if b.err != nil {
			return n - remain, b.readErr()
		}
	}
}

func (b *bufferedReader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}

	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			// large read, empty buffer
			// read directly into p to avoid extra copy
			n, b.err = b.rd.Read(p)
			if n < 0 {
				return n, fmt.Errorf("arrow/bufferedreader: %w", bufio.ErrNegativeCount)
			}
			return n, b.readErr()
		}

		// one read
		// don't use b.fill
		b.r, b.w = 0, 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			return n, fmt.Errorf("arrow/bufferedreader: %w", bufio.ErrNegativeCount)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}
