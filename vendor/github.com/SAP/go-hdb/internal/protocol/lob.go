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
var lobChunkSize int32 = 256 //TODO: check size

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

// LobReadDescr is the package internal representation of a lob field to be read from database.
type LobReadDescr struct {
	col int
	fn  func() error
	w   lobWriter
}

// SetWriter sets the io.Writer destination for a lob field to be read from database.
func (d *LobReadDescr) SetWriter(w io.Writer) error {
	if err := d.w.setWriter(w); err != nil {
		return err
	}
	if d.fn != nil {
		return d.fn()
	}
	return nil
}

// LobWriteDescr is the package internal representation of a lob field to be written to database.
type LobWriteDescr struct {
	r io.Reader
}

// SetReader sets the io.Reader source for a lob field to be written to database.
func (d *LobWriteDescr) SetReader(r io.Reader) {
	d.r = r
}

type locatorID uint64 // byte[locatorIdSize]

// write lob reply
type writeLobReply struct {
	ids    []locatorID
	numArg int
}

func newWriteLobReply() *writeLobReply {
	return &writeLobReply{
		ids: make([]locatorID, 0),
	}
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
	if cap(r.ids) < r.numArg {
		r.ids = make([]locatorID, r.numArg)
	} else {
		r.ids = r.ids[:r.numArg]
	}

	for i := 0; i < r.numArg; i++ {
		if id, err := rd.ReadUint64(); err == nil {
			r.ids[i] = locatorID(id)
		} else {
			return err
		}
	}

	return nil
}

//write lob request
type writeLobRequest struct {
	readers []lobReader
}

func newWriteLobRequest(readers []lobReader) *writeLobRequest {
	return &writeLobRequest{
		readers: readers,
	}
}

func (r *writeLobRequest) kind() partKind {
	return pkWriteLobRequest
}

func (r *writeLobRequest) size() (int, error) {

	// TODO: check size limit

	size := 0
	for _, reader := range r.readers {
		if reader.done() {
			continue
		}

		if err := reader.fill(); err != nil {
			return 0, err
		}
		size += writeLobRequestHeaderSize
		size += reader.size()
	}
	return size, nil
}

func (r *writeLobRequest) numArg() int {
	n := 0
	for _, reader := range r.readers {
		if !reader.done() {
			n++
		}
	}
	return n
}

func (r *writeLobRequest) write(wr *bufio.Writer) error {
	for _, reader := range r.readers {
		if !reader.done() {

			if err := wr.WriteUint64(uint64(reader.id())); err != nil {
				return err
			}

			opt := int8(0x02) // data included
			if reader.eof() {
				opt |= 0x04 // last data
			}

			if err := wr.WriteInt8(opt); err != nil {
				return err
			}

			if err := wr.WriteInt64(-1); err != nil { //offset (-1 := append)
				return err
			}

			if err := wr.WriteInt32(int32(reader.size())); err != nil { // size
				return err
			}

			if _, err := wr.Write(reader.bytes()); err != nil {
				return err
			}
		}
	}
	return nil
}

//read lob request
type readLobRequest struct {
	writers []lobWriter
}

func (r *readLobRequest) numWriter() int {
	n := 0
	for _, writer := range r.writers {
		if !writer.eof() {
			n++
		}
	}
	return n
}

func (r *readLobRequest) kind() partKind {
	return pkReadLobRequest
}

func (r *readLobRequest) size() (int, error) {
	return r.numWriter() * readLobRequestSize, nil
}

func (r *readLobRequest) numArg() int {
	return r.numWriter()
}

func (r *readLobRequest) write(wr *bufio.Writer) error {
	for _, writer := range r.writers {
		if writer.eof() {
			continue
		}

		if err := wr.WriteUint64(uint64(writer.id())); err != nil {
			return err
		}

		readOfs, readLen := writer.readOfsLen()

		if err := wr.WriteInt64(readOfs + 1); err != nil { //1-based
			return err
		}

		if err := wr.WriteInt32(readLen); err != nil {
			return err
		}

		if err := wr.WriteZeroes(4); err != nil {
			return err
		}
	}
	return nil
}

// read lob reply
// - seems like readLobreply gives only an result for one lob - even if more then one is requested
type readLobReply struct {
	writers []lobWriter
	numArg  int
}

func (r *readLobReply) kind() partKind {
	return pkReadLobReply
}

func (r *readLobReply) setNumArg(numArg int) {
	r.numArg = numArg
}

func (r *readLobReply) read(rd *bufio.Reader) error {
	for i := 0; i < r.numArg; i++ {

		id, err := rd.ReadUint64()
		if err != nil {
			return err
		}

		var writer lobWriter
		for _, writer = range r.writers {
			if writer.id() == locatorID(id) {
				break // writer found
			}
		}
		if writer == nil {
			return fmt.Errorf("internal error: no lob writer found for id %d", id)
		}

		opt, err := rd.ReadInt8()
		if err != nil {
			return err
		}

		chunkLen, err := rd.ReadInt32()
		if err != nil {
			return err
		}

		if err := rd.Skip(3); err != nil {
			return err
		}

		eof := (lobOptions(opt) & loLastdata) != 0

		if err := writer.write(rd, int(chunkLen), eof); err != nil {
			return err
		}
	}
	return nil
}

// lobWriter reads lob chunks and writes them into lob field.
type lobWriter interface {
	id() locatorID
	setWriter(w io.Writer) error
	write(rd *bufio.Reader, size int, eof bool) error
	readOfsLen() (int64, int32)
	eof() bool
}

// baseLobWriter is a reuse struct for binary and char lob writers.
type baseLobWriter struct {
	_id     locatorID
	charLen int64
	byteLen int64

	readOfs int64
	_eof    bool

	ofs int

	wr io.Writer

	_flush func() error

	b []byte
}

func (l *baseLobWriter) id() locatorID {
	return l._id
}

func (l *baseLobWriter) eof() bool {
	return l._eof
}

func (l *baseLobWriter) setWriter(wr io.Writer) error {
	l.wr = wr
	return l._flush()
}

func (l *baseLobWriter) write(rd *bufio.Reader, size int, eof bool) error {
	l._eof = eof // store eof

	if size == 0 {
		return nil
	}

	l.b = resizeBuffer(l.b, size+l.ofs)
	if err := rd.ReadFull(l.b[l.ofs:]); err != nil {
		return err
	}
	if l.wr != nil {
		return l._flush()
	}
	return nil
}

func (l *baseLobWriter) readOfsLen() (int64, int32) {
	readLen := l.charLen - l.readOfs
	if readLen > int64(math.MaxInt32) || readLen > int64(lobChunkSize) {
		return l.readOfs, lobChunkSize
	}
	return l.readOfs, int32(readLen)
}

// binaryLobWriter (byte based lobs).
type binaryLobWriter struct {
	*baseLobWriter
}

func newBinaryLobWriter(id locatorID, charLen, byteLen int64) *binaryLobWriter {
	l := &binaryLobWriter{
		baseLobWriter: &baseLobWriter{_id: id, charLen: charLen, byteLen: byteLen},
	}
	l._flush = l.flush
	return l
}

func (l *binaryLobWriter) flush() error {
	if _, err := l.wr.Write(l.b); err != nil {
		return err
	}
	l.readOfs += int64(len(l.b))
	return nil
}

type charLobWriter struct {
	*baseLobWriter
}

func newCharLobWriter(id locatorID, charLen, byteLen int64) *charLobWriter {
	l := &charLobWriter{
		baseLobWriter: &baseLobWriter{_id: id, charLen: charLen, byteLen: byteLen},
	}
	l._flush = l.flush
	return l
}

func (l *charLobWriter) flush() error {
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
func (l *charLobWriter) runeCount(b []byte) int {
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

// lobWriter reads field lob data chunks.
type lobReader interface {
	id() locatorID
	fill() error
	size() int
	bytes() []byte
	eof() bool
	done() bool
}

// baseLobWriter is a reuse struct for binary and char lob writers.
type baseLobReader struct {
	r     io.Reader
	_id   locatorID
	_size int
	_eof  bool
	_done bool
	b     []byte
}

func (l *baseLobReader) id() locatorID {
	return l._id
}

func (l *baseLobReader) eof() bool {
	return l._eof
}

func (l *baseLobReader) done() bool {
	return l._done
}

func (l *baseLobReader) size() int {
	return l._size
}

func (l *baseLobReader) bytes() []byte {
	if l._eof {
		l._done = true
	}
	return l.b[:l._size]
}

// binaryLobReader (byte based lobs).
type binaryLobReader struct {
	*baseLobReader
}

func newBinaryLobReader(r io.Reader, id locatorID) *binaryLobReader {
	return &binaryLobReader{
		baseLobReader: &baseLobReader{r: r, _id: id},
	}
}

func (l *binaryLobReader) fill() error {
	if l._eof {
		return fmt.Errorf("locator id %d eof error", l._id)
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

// charLobReader (character based lobs - cesu8).
type charLobReader struct {
	*baseLobReader
	c   []byte
	ofs int
}

func newCharLobReader(r io.Reader, id locatorID) *charLobReader {
	return &charLobReader{
		baseLobReader: &baseLobReader{r: r, _id: id},
	}
}

func (l *charLobReader) fill() error {
	if l._eof {
		return fmt.Errorf("locator id %d eof error", l._id)
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

// helper
func resizeBuffer(b1 []byte, size int) []byte {
	if b1 == nil || cap(b1) < size {
		b2 := make([]byte, size)
		copy(b2, b1) // !!!
		return b2
	}
	return b1[:size]
}
