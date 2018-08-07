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

package protocol

import (
	"fmt"
	"io"
	"math"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/SAP/go-hdb/internal/bufio"
	"github.com/SAP/go-hdb/internal/unicode"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

const (
	locatorIDSize             = 8
	writeLobRequestHeaderSize = 21
	readLobRequestSize        = 24
)

// variable (unit testing)
//var lobChunkSize = 1 << 14 //TODO: check size
var lobChunkSize int32 = 4096 //TODO: check size

//lob options
type lobOptions int8

const (
	loNullindicator lobOptions = 0x01
	loDataincluded  lobOptions = 0x02
	loLastdata      lobOptions = 0x04
)

var lobOptionsText = map[lobOptions]string{
	loNullindicator: "null indicator",
	loDataincluded:  "data included",
	loLastdata:      "last data",
}

func (k lobOptions) String() string {
	t := make([]string, 0, len(lobOptionsText))

	for option, text := range lobOptionsText {
		if (k & option) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

type locatorID uint64 // byte[locatorIdSize]

// write lob reply
type writeLobReply struct {
	ids    []locatorID
	numArg int
}

func (r *writeLobReply) String() string {
	return fmt.Sprintf("write lob reply: %v", r.ids)
}

func (r *writeLobReply) kind() partKind {
	return pkWriteLobReply
}

func (r *writeLobReply) setNumArg(numArg int) {
	r.numArg = numArg
}

func (r *writeLobReply) read(rd *bufio.Reader) error {

	//resize ids
	if r.ids == nil || cap(r.ids) < r.numArg {
		r.ids = make([]locatorID, r.numArg)
	} else {
		r.ids = r.ids[:r.numArg]
	}

	for i := 0; i < r.numArg; i++ {
		r.ids[i] = locatorID(rd.ReadUint64())
	}

	return rd.GetError()
}

//write lob request
type writeLobRequest struct {
	lobPrmFields []*ParameterField
}

func (r *writeLobRequest) kind() partKind {
	return pkWriteLobRequest
}

func (r *writeLobRequest) size() (int, error) {

	// TODO: check size limit

	size := 0
	for _, prmField := range r.lobPrmFields {
		cr := prmField.chunkReader
		if cr.done() {
			continue
		}

		if err := cr.fill(); err != nil {
			return 0, err
		}
		size += writeLobRequestHeaderSize
		size += cr.size()
	}
	return size, nil
}

func (r *writeLobRequest) numArg() int {
	n := 0
	for _, prmField := range r.lobPrmFields {
		cr := prmField.chunkReader
		if !cr.done() {
			n++
		}
	}
	return n
}

func (r *writeLobRequest) write(wr *bufio.Writer) error {
	for _, prmField := range r.lobPrmFields {
		cr := prmField.chunkReader
		if !cr.done() {

			wr.WriteUint64(uint64(prmField.lobLocatorID))

			opt := int8(0x02) // data included
			if cr.eof() {
				opt |= 0x04 // last data
			}

			wr.WriteInt8(opt)
			wr.WriteInt64(-1)               //offset (-1 := append)
			wr.WriteInt32(int32(cr.size())) // size
			wr.Write(cr.bytes())
		}
	}
	return nil
}

//read lob request
type readLobRequest struct {
	w lobChunkWriter
}

func (r *readLobRequest) kind() partKind {
	return pkReadLobRequest
}

func (r *readLobRequest) size() (int, error) {
	return readLobRequestSize, nil
}

func (r *readLobRequest) numArg() int {
	return 1
}

func (r *readLobRequest) write(wr *bufio.Writer) error {
	wr.WriteUint64(uint64(r.w.id()))

	readOfs, readLen := r.w.readOfsLen()

	wr.WriteInt64(readOfs + 1) //1-based
	wr.WriteInt32(readLen)
	wr.WriteZeroes(4)

	return nil
}

// read lob reply
// - seems like readLobreply gives only an result for one lob - even if more then one is requested
// --> read single lobs
type readLobReply struct {
	w lobChunkWriter
}

func (r *readLobReply) kind() partKind {
	return pkReadLobReply
}

func (r *readLobReply) setNumArg(numArg int) {
	if numArg != 1 {
		panic("numArg == 1 expected")
	}
}

func (r *readLobReply) read(rd *bufio.Reader) error {
	id := rd.ReadUint64()

	if r.w.id() != locatorID(id) {
		return fmt.Errorf("internal error: invalid lob locator %d - expected %d", id, r.w.id())
	}

	opt := rd.ReadInt8()
	chunkLen := rd.ReadInt32()
	rd.Skip(3)
	eof := (lobOptions(opt) & loLastdata) != 0

	if err := r.w.write(rd, int(chunkLen), eof); err != nil {
		return err
	}

	return rd.GetError()
}

// lobChunkReader reads lob field io.Reader in chunks for writing to db.
type lobChunkReader interface {
	fill() error
	size() int
	bytes() []byte
	eof() bool
	done() bool
}

func newLobChunkReader(isCharBased bool, r io.Reader) lobChunkReader {
	if isCharBased {
		return &charLobChunkReader{r: r}
	}
	return &binaryLobChunkReader{r: r}
}

// binaryLobChunkReader (byte based chunks).
type binaryLobChunkReader struct {
	r     io.Reader
	_size int
	_eof  bool
	_done bool
	b     []byte
}

func (l *binaryLobChunkReader) eof() bool  { return l._eof }
func (l *binaryLobChunkReader) done() bool { return l._done }
func (l *binaryLobChunkReader) size() int  { return l._size }

func (l *binaryLobChunkReader) bytes() []byte {
	l._done = l._eof
	return l.b[:l._size]
}

func (l *binaryLobChunkReader) fill() error {
	if l._eof {
		return io.EOF
	}

	var err error

	l.b = resizeBuffer(l.b, int(lobChunkSize))
	l._size, err = l.r.Read(l.b)
	if err != nil && err != io.EOF {
		return err
	}
	l._eof = err == io.EOF
	return nil
}

// charLobChunkReader (cesu8 character based chunks).
type charLobChunkReader struct {
	r     io.Reader
	_size int
	_eof  bool
	_done bool
	b     []byte
	c     []byte
	ofs   int
}

func (l *charLobChunkReader) eof() bool  { return l._eof }
func (l *charLobChunkReader) done() bool { return l._done }
func (l *charLobChunkReader) size() int  { return l._size }

func (l *charLobChunkReader) bytes() []byte {
	l._done = l._eof
	return l.b[:l._size]
}

func (l *charLobChunkReader) fill() error {
	if l._eof {
		return io.EOF
	}

	l.c = resizeBuffer(l.c, int(lobChunkSize)+l.ofs)
	n, err := l.r.Read(l.c[l.ofs:])
	size := n + l.ofs

	if err != nil && err != io.EOF {
		return err
	}
	l._eof = err == io.EOF
	if l._eof && size == 0 {
		l._size = 0
		return nil
	}

	l.b = resizeBuffer(l.b, cesu8.Size(l.c[:size])) // last rune might be incomplete, so size is one greater than needed
	nDst, nSrc, err := unicode.Utf8ToCesu8Transformer.Transform(l.b, l.c[:size], l._eof)
	if err != nil && err != transform.ErrShortSrc {
		return err
	}

	if l._eof && err == transform.ErrShortSrc {
		return unicode.ErrInvalidUtf8
	}

	l._size = nDst
	l.ofs = size - nSrc

	if l.ofs > 0 {
		copy(l.c, l.c[nSrc:size]) // copy rest to buffer beginn
	}
	return nil
}

// lobChunkWriter reads db lob chunks and writes them into lob field io.Writer.
type lobChunkWriter interface {
	SetWriter(w io.Writer) error // gets called by driver.Lob.Scan

	id() locatorID
	write(rd *bufio.Reader, size int, eof bool) error
	readOfsLen() (int64, int32)
	eof() bool
}

func newLobChunkWriter(isCharBased bool, s *Session, id locatorID, charLen, byteLen int64) lobChunkWriter {
	if isCharBased {
		return &charLobChunkWriter{s: s, _id: id, charLen: charLen, byteLen: byteLen}
	}
	return &binaryLobChunkWriter{s: s, _id: id, charLen: charLen, byteLen: byteLen}
}

// binaryLobChunkWriter (byte based lobs).
type binaryLobChunkWriter struct {
	s *Session

	_id     locatorID
	charLen int64
	byteLen int64

	readOfs int64
	_eof    bool

	ofs int

	wr io.Writer

	b []byte
}

func (l *binaryLobChunkWriter) id() locatorID { return l._id }
func (l *binaryLobChunkWriter) eof() bool     { return l._eof }

func (l *binaryLobChunkWriter) SetWriter(wr io.Writer) error {
	l.wr = wr
	if err := l.flush(); err != nil {
		return err
	}
	return l.s.readLobStream(l)
}

func (l *binaryLobChunkWriter) write(rd *bufio.Reader, size int, eof bool) error {
	l._eof = eof // store eof

	if size == 0 {
		return nil
	}

	l.b = resizeBuffer(l.b, size+l.ofs)
	rd.ReadFull(l.b[l.ofs:])
	if l.wr != nil {
		return l.flush()
	}
	return nil
}

func (l *binaryLobChunkWriter) readOfsLen() (int64, int32) {
	readLen := l.charLen - l.readOfs
	if readLen > int64(math.MaxInt32) || readLen > int64(lobChunkSize) {
		return l.readOfs, lobChunkSize
	}
	return l.readOfs, int32(readLen)
}

func (l *binaryLobChunkWriter) flush() error {
	if _, err := l.wr.Write(l.b); err != nil {
		return err
	}
	l.readOfs += int64(len(l.b))
	return nil
}

type charLobChunkWriter struct {
	s *Session

	_id     locatorID
	charLen int64
	byteLen int64

	readOfs int64
	_eof    bool

	ofs int

	wr io.Writer

	b []byte
}

func (l *charLobChunkWriter) id() locatorID { return l._id }
func (l *charLobChunkWriter) eof() bool     { return l._eof }

func (l *charLobChunkWriter) SetWriter(wr io.Writer) error {
	l.wr = wr
	if err := l.flush(); err != nil {
		return err
	}
	return l.s.readLobStream(l)
}

func (l *charLobChunkWriter) write(rd *bufio.Reader, size int, eof bool) error {
	l._eof = eof // store eof

	if size == 0 {
		return nil
	}

	l.b = resizeBuffer(l.b, size+l.ofs)
	rd.ReadFull(l.b[l.ofs:])
	if l.wr != nil {
		return l.flush()
	}
	return nil
}

func (l *charLobChunkWriter) readOfsLen() (int64, int32) {
	readLen := l.charLen - l.readOfs
	if readLen > int64(math.MaxInt32) || readLen > int64(lobChunkSize) {
		return l.readOfs, lobChunkSize
	}
	return l.readOfs, int32(readLen)
}

func (l *charLobChunkWriter) flush() error {
	nDst, nSrc, err := unicode.Cesu8ToUtf8Transformer.Transform(l.b, l.b, true) // inline cesu8 to utf8 transformation
	if err != nil && err != transform.ErrShortSrc {
		return err
	}
	if _, err := l.wr.Write(l.b[:nDst]); err != nil {
		return err
	}
	l.ofs = len(l.b) - nSrc
	if l.ofs != 0 && l.ofs != cesu8.CESUMax/2 { // assert remaining bytes
		return unicode.ErrInvalidCesu8
	}
	l.readOfs += int64(l.runeCount(l.b[:nDst]))
	if l.ofs != 0 {
		l.readOfs++                   // add half encoding
		copy(l.b, l.b[nSrc:len(l.b)]) // move half encoding to buffer begin
	}
	return nil
}

// Caution: hdb counts 4 byte utf-8 encodings (cesu-8 6 bytes) as 2 (3 byte) chars
func (l *charLobChunkWriter) runeCount(b []byte) int {
	numChars := 0
	for len(b) > 0 {
		_, size := utf8.DecodeRune(b)
		b = b[size:]
		numChars++
		if size == utf8.UTFMax {
			numChars++
		}
	}
	return numChars
}

// helper
func resizeBuffer(b1 []byte, size int) []byte {
	if b1 == nil || cap(b1) < size {
		b2 := make([]byte, size)
		copy(b2, b1) // !!!
		return b2
	}
	return b1[:size]
}
