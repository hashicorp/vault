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

// Package bufio implements buffered I/O for database read and writes on basis of the standard Go bufio package.
package bufio

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"

	"github.com/SAP/go-hdb/internal/unicode"
	"golang.org/x/text/transform"
)

const (
	bufferSize = 128
)

// Reader is a bufio.Reader extended by methods needed for hdb protocol.
type Reader struct {
	*bufio.Reader
	b  []byte // scratch buffer (min 8 Bytes)
	tr transform.Transformer
}

// NewReader creates a new Reader instance.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Reader: bufio.NewReader(r),
		b:      make([]byte, bufferSize),
		tr:     unicode.Cesu8ToUtf8Transformer,
	}
}

// NewReaderSize creates a new Reader instance with given size for bufio.Reader.
func NewReaderSize(r io.Reader, size int) *Reader {
	return &Reader{
		Reader: bufio.NewReaderSize(r, size),
		b:      make([]byte, bufferSize),
		tr:     unicode.Cesu8ToUtf8Transformer,
	}
}

// Skip skips cnt bytes from reading.
func (r *Reader) Skip(cnt int) error {
	for i := 0; i < cnt; {
		j := cnt - i
		if j > len(r.b) {
			j = len(r.b)
		}
		n, err := io.ReadFull(r.Reader, r.b[:j])
		i += n
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadFull implements io.ReadFull on Reader.
func (r *Reader) ReadFull(p []byte) error {
	_, err := io.ReadFull(r.Reader, p)
	return err
}

// ReadBool reads and returns a boolean.
func (r *Reader) ReadBool() (bool, error) {
	c, err := r.Reader.ReadByte()
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

// ReadInt8 reads and returns an int8.
func (r *Reader) ReadInt8() (int8, error) {
	c, err := r.Reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return int8(c), nil
}

// ReadInt16 reads and returns an int16.
func (r *Reader) ReadInt16() (int16, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:2]); err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(r.b[:2])), nil
}

// ReadUint16 reads and returns an uint16.
func (r *Reader) ReadUint16() (uint16, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:2]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(r.b[:2]), nil
}

// ReadInt32 reads and returns an int32.
func (r *Reader) ReadInt32() (int32, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:4]); err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(r.b[:4])), nil
}

// ReadUint32 reads and returns an uint32.
func (r *Reader) ReadUint32() (uint32, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:4]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(r.b[:4]), nil
}

// ReadInt64 reads and returns an int64.
func (r *Reader) ReadInt64() (int64, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:8]); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(r.b[:8])), nil
}

// ReadUint64 reads and returns an uint64.
func (r *Reader) ReadUint64() (uint64, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:8]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(r.b[:8]), nil
}

// ReadFloat32 reads and returns a float32.
func (r *Reader) ReadFloat32() (float32, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:4]); err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(r.b[:4])
	return math.Float32frombits(bits), nil
}

// ReadFloat64 reads and returns a float64.
func (r *Reader) ReadFloat64() (float64, error) {
	if _, err := io.ReadFull(r.Reader, r.b[:8]); err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(r.b[:8])
	return math.Float64frombits(bits), nil
}

// ReadCesu8 reads a size CESU-8 encoded byte sequence and returns an UTF-8 byte slice.
func (r *Reader) ReadCesu8(size int) ([]byte, error) {
	p := make([]byte, size)
	if _, err := io.ReadFull(r.Reader, p); err != nil {
		return nil, err
	}
	r.tr.Reset()
	n, _, err := r.tr.Transform(p, p, true) // inplace transformation
	if err != nil {
		return nil, err
	}
	return p[:n], nil
}

// Writer is a bufio.Writer extended by methods needed for hdb protocol.
type Writer struct {
	*bufio.Writer
	b  []byte // // scratch buffer (min 8 Bytes)
	tr transform.Transformer
}

// NewWriter creates a new Writer instance.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: bufio.NewWriter(w),
		b:      make([]byte, bufferSize),
		tr:     unicode.Utf8ToCesu8Transformer,
	}
}

// NewWriterSize creates a new Writer instance with given size for bufio.Writer.
func NewWriterSize(w io.Writer, size int) *Writer {
	return &Writer{
		Writer: bufio.NewWriterSize(w, size),
		b:      make([]byte, bufferSize),
		tr:     unicode.Utf8ToCesu8Transformer,
	}
}

// WriteZeroes writes cnt zero byte values.
func (w *Writer) WriteZeroes(cnt int) error {
	// zero out scratch area
	l := cnt
	if l > len(w.b) {
		l = len(w.b)
	}
	for i := 0; i < l; i++ {
		w.b[i] = 0
	}

	for i := 0; i < cnt; {
		j := cnt - i
		if j > len(w.b) {
			j = len(w.b)
		}
		n, err := w.Writer.Write(w.b[:j])
		i += n
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteBool writes a boolean.
func (w *Writer) WriteBool(v bool) error {
	if v {
		return w.Writer.WriteByte(1)
	}
	return w.Writer.WriteByte(0)
}

// WriteInt8 writes an int8.
func (w *Writer) WriteInt8(i int8) error {
	return w.Writer.WriteByte(byte(i))
}

// WriteInt16 writes an int16.
func (w *Writer) WriteInt16(i int16) error {
	binary.LittleEndian.PutUint16(w.b[:2], uint16(i))
	_, err := w.Writer.Write(w.b[:2])
	return err
}

// WriteUint16 writes an uint16.
func (w *Writer) WriteUint16(i uint16) error {
	binary.LittleEndian.PutUint16(w.b[:2], i)
	_, err := w.Writer.Write(w.b[:2])
	return err
}

// WriteInt32 writes an int32.
func (w *Writer) WriteInt32(i int32) error {
	binary.LittleEndian.PutUint32(w.b[:4], uint32(i))
	_, err := w.Writer.Write(w.b[:4])
	return err
}

// WriteUint32 writes an uint32.
func (w *Writer) WriteUint32(i uint32) error {
	binary.LittleEndian.PutUint32(w.b[:4], i)
	_, err := w.Writer.Write(w.b[:4])
	return err
}

// WriteInt64 writes an int64.
func (w *Writer) WriteInt64(i int64) error {
	binary.LittleEndian.PutUint64(w.b[:8], uint64(i))
	_, err := w.Writer.Write(w.b[:8])
	return err
}

// WriteUint64 writes an uint64.
func (w *Writer) WriteUint64(i uint64) error {
	binary.LittleEndian.PutUint64(w.b[:8], i)
	_, err := w.Writer.Write(w.b[:8])
	return err
}

// WriteFloat32 writes a float32.
func (w *Writer) WriteFloat32(f float32) error {
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(w.b[:4], bits)
	_, err := w.Writer.Write(w.b[:4])
	return err
}

// WriteFloat64 writes a float64.
func (w *Writer) WriteFloat64(f float64) error {
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(w.b[:8], bits)
	_, err := w.Writer.Write(w.b[:8])
	return err
}

// WriteCesu8 writes an UTF-8 byte slice as CESU-8 and returns the CESU-8 bytes written.
func (w *Writer) WriteCesu8(p []byte) (int, error) {
	w.tr.Reset()
	cnt := 0
	i := 0
	for i < len(p) {
		m, n, err := w.tr.Transform(w.b, p[i:], true)
		if err != nil && err != transform.ErrShortDst {
			return cnt, err
		}
		if m == 0 {
			return cnt, transform.ErrShortDst
		}
		o, err := w.Writer.Write(w.b[:m])
		cnt += o
		if err != nil {
			return cnt, err
		}
		i += n
	}
	return cnt, nil
}

// WriteStringCesu8 is like WriteCesu8 with an UTF-8 string as parameter.
func (w *Writer) WriteStringCesu8(s string) (int, error) {
	return w.WriteCesu8([]byte(s))
}
