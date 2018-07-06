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

// Reader is a bufio.Reader extended by methods needed for hdb protocol.
type Reader struct {
	rd  *bufio.Reader
	err error
	b   [8]byte // scratch buffer (8 Bytes)
	tr  transform.Transformer
	cnt int
}

// NewReader creates a new Reader instance.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(r),
		tr: unicode.Cesu8ToUtf8Transformer,
	}
}

// NewReaderSize creates a new Reader instance with given size for bufio.Reader.
func NewReaderSize(r io.Reader, size int) *Reader {
	return &Reader{
		rd: bufio.NewReaderSize(r, size),
		tr: unicode.Cesu8ToUtf8Transformer,
	}
}

// ResetCnt resets the byte read counter.
func (r *Reader) ResetCnt() {
	r.cnt = 0
}

// Cnt returns the value of the byte read counter.
func (r *Reader) Cnt() int {
	return r.cnt
}

// GetError returns reader error.
func (r *Reader) GetError() error {
	err := r.err
	r.err = nil
	return err
}

// Skip skips cnt bytes from reading.
func (r *Reader) Skip(cnt int) {
	if r.err != nil {
		return
	}
	var n int
	n, r.err = r.rd.Discard(cnt)
	r.cnt += n
}

// ReadB reads and returns a byte.
func (r *Reader) ReadB() byte { // ReadB as sig differs from ReadByte (vet issues)
	if r.err != nil {
		return 0
	}
	var b byte
	b, r.err = r.rd.ReadByte()
	r.cnt++
	return b
}

// ReadFull implements io.ReadFull on Reader.
func (r *Reader) ReadFull(p []byte) {
	if r.err != nil {
		return
	}
	var n int
	n, r.err = io.ReadFull(r.rd, p)
	r.cnt += n
}

// ReadBool reads and returns a boolean.
func (r *Reader) ReadBool() bool {
	if r.err != nil {
		return false
	}
	return !(r.ReadB() == 0)
}

// ReadInt8 reads and returns an int8.
func (r *Reader) ReadInt8() int8 {
	return int8(r.ReadB())
}

// ReadInt16 reads and returns an int16.
func (r *Reader) ReadInt16() int16 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:2])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return int16(binary.LittleEndian.Uint16(r.b[:2]))
}

// ReadUint16 reads and returns an uint16.
func (r *Reader) ReadUint16() uint16 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:2])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint16(r.b[:2])
}

// ReadInt32 reads and returns an int32.
func (r *Reader) ReadInt32() int32 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:4])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(r.b[:4]))
}

// ReadUint32 reads and returns an uint32.
func (r *Reader) ReadUint32() uint32 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:4])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(r.b[:4])
}

// ReadInt64 reads and returns an int64.
func (r *Reader) ReadInt64() int64 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:8])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(r.b[:8]))
}

// ReadUint64 reads and returns an uint64.
func (r *Reader) ReadUint64() uint64 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:8])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(r.b[:8])
}

// ReadFloat32 reads and returns a float32.
func (r *Reader) ReadFloat32() float32 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:4])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint32(r.b[:4])
	return math.Float32frombits(bits)
}

// ReadFloat64 reads and returns a float64.
func (r *Reader) ReadFloat64() float64 {
	if r.err != nil {
		return 0
	}
	var n int
	n, r.err = io.ReadFull(r.rd, r.b[:8])
	r.cnt += n
	if r.err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint64(r.b[:8])
	return math.Float64frombits(bits)
}

// ReadCesu8 reads a size CESU-8 encoded byte sequence and returns an UTF-8 byte slice.
func (r *Reader) ReadCesu8(size int) []byte {
	if r.err != nil {
		return nil
	}
	p := make([]byte, size)
	var n int
	n, r.err = io.ReadFull(r.rd, p)
	r.cnt += n
	if r.err != nil {
		return nil
	}
	r.tr.Reset()
	if n, _, r.err = r.tr.Transform(p, p, true); r.err != nil { // inplace transformation
		return nil
	}
	return p[:n]
}

const writerBufferSize = 4096

// Writer is a bufio.Writer extended by methods needed for hdb protocol.
type Writer struct {
	wr  *bufio.Writer
	err error
	b   []byte // scratch buffer (min 8 Bytes)
	tr  transform.Transformer
}

// NewWriter creates a new Writer instance.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		wr: bufio.NewWriter(w),
		b:  make([]byte, writerBufferSize),
		tr: unicode.Utf8ToCesu8Transformer,
	}
}

// NewWriterSize creates a new Writer instance with given size for bufio.Writer.
func NewWriterSize(w io.Writer, size int) *Writer {
	return &Writer{
		wr: bufio.NewWriterSize(w, size),
		b:  make([]byte, writerBufferSize),
		tr: unicode.Utf8ToCesu8Transformer,
	}
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	if w.err != nil {
		return w.err
	}
	return w.wr.Flush()
}

// WriteZeroes writes cnt zero byte values.
func (w *Writer) WriteZeroes(cnt int) {
	if w.err != nil {
		return
	}

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
		n, _ := w.wr.Write(w.b[:j])
		i += n
	}
}

// Write writes the contents of p.
func (w *Writer) Write(p []byte) {
	if w.err != nil {
		return
	}
	w.wr.Write(p)
}

// WriteB writes a byte.
func (w *Writer) WriteB(b byte) { // WriteB as sig differs from WriteByte (vet issues)
	if w.err != nil {
		return
	}
	w.wr.WriteByte(b)
}

// WriteBool writes a boolean.
func (w *Writer) WriteBool(v bool) {
	if w.err != nil {
		return
	}
	if v {
		w.wr.WriteByte(1)
	} else {
		w.wr.WriteByte(0)
	}
}

// WriteInt8 writes an int8.
func (w *Writer) WriteInt8(i int8) {
	if w.err != nil {
		return
	}
	w.wr.WriteByte(byte(i))
}

// WriteInt16 writes an int16.
func (w *Writer) WriteInt16(i int16) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint16(w.b[:2], uint16(i))
	w.wr.Write(w.b[:2])
}

// WriteUint16 writes an uint16.
func (w *Writer) WriteUint16(i uint16) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint16(w.b[:2], i)
	w.wr.Write(w.b[:2])
}

// WriteInt32 writes an int32.
func (w *Writer) WriteInt32(i int32) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint32(w.b[:4], uint32(i))
	w.wr.Write(w.b[:4])
}

// WriteUint32 writes an uint32.
func (w *Writer) WriteUint32(i uint32) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint32(w.b[:4], i)
	w.wr.Write(w.b[:4])
}

// WriteInt64 writes an int64.
func (w *Writer) WriteInt64(i int64) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint64(w.b[:8], uint64(i))
	w.wr.Write(w.b[:8])
}

// WriteUint64 writes an uint64.
func (w *Writer) WriteUint64(i uint64) {
	if w.err != nil {
		return
	}
	binary.LittleEndian.PutUint64(w.b[:8], i)
	w.wr.Write(w.b[:8])
}

// WriteFloat32 writes a float32.
func (w *Writer) WriteFloat32(f float32) {
	if w.err != nil {
		return
	}
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(w.b[:4], bits)
	w.wr.Write(w.b[:4])
}

// WriteFloat64 writes a float64.
func (w *Writer) WriteFloat64(f float64) {
	if w.err != nil {
		return
	}
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(w.b[:8], bits)
	w.wr.Write(w.b[:8])
}

// WriteString writes a string.
func (w *Writer) WriteString(s string) {
	if w.err != nil {
		return
	}
	w.wr.WriteString(s)
}

// WriteCesu8 writes an UTF-8 byte slice as CESU-8 and returns the CESU-8 bytes written.
func (w *Writer) WriteCesu8(p []byte) int {
	if w.err != nil {
		return 0
	}
	w.tr.Reset()
	cnt := 0
	i := 0
	for i < len(p) {
		m, n, err := w.tr.Transform(w.b, p[i:], true)
		if err != nil && err != transform.ErrShortDst {
			w.err = err
			return cnt
		}
		if m == 0 {
			w.err = transform.ErrShortDst
			return cnt
		}
		o, _ := w.wr.Write(w.b[:m])
		cnt += o
		i += n
	}
	return cnt
}

// WriteStringCesu8 is like WriteCesu8 with an UTF-8 string as parameter.
func (w *Writer) WriteStringCesu8(s string) int {
	return w.WriteCesu8([]byte(s))
}
