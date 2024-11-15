package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"sync"
	"unicode/utf8"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
	"github.com/SAP/go-hdb/driver/internal/unsafe"
	"golang.org/x/text/transform"
)

const (
	writeLobRequestSize = 21
)

// lobOptions represents a lob option set.
type lobOptions int8

const (
	loNullindicator lobOptions = 0x01
	loDataincluded  lobOptions = 0x02
	loLastdata      lobOptions = 0x04
)

const (
	loNullindicatorText = "null indicator"
	loDataincludedText  = "data included"
	loLastdataText      = "last data"
)

func (o lobOptions) String() string {
	var s []string
	if o&loNullindicator != 0 {
		s = append(s, loNullindicatorText)
	}
	if o&loDataincluded != 0 {
		s = append(s, loDataincludedText)
	}
	if o&loLastdata != 0 {
		s = append(s, loLastdataText)
	}
	return fmt.Sprintf("%v", s)
}

// IsLastData return true if the last data package was read, false otherwise.
func (o lobOptions) isLastData() bool { return (o & loLastdata) != 0 }
func (o lobOptions) isNull() bool     { return (o & loNullindicator) != 0 }

// lob typecode.
type lobTypecode int8

const (
	ltcUndefined lobTypecode = 0
	ltcBlob      lobTypecode = 1
	ltcClob      lobTypecode = 2
	ltcNclob     lobTypecode = 3
)

// not used
// type lobFlags bool

// func (f lobFlags) String() string { return fmt.Sprintf("%t", f) }
// func (f *lobFlags) decode(dec *encoding.Decoder, ph *partHeader) error {
// 	*f = lobFlags(dec.Bool())
// 	return dec.Error()
// }
// func (f lobFlags) encode(enc *encoding.Encoder) error { enc.Bool(bool(f)); return nil }

// LobScanner is the interface wrapping the Scan method for Lob reading.
type LobScanner interface {
	Scan(w io.Writer) error
}

var _ LobScanner = (*lobOutDescr)(nil)

// LobInDescr represents a lob input descriptor.
type LobInDescr struct {
	rd  io.Reader
	opt lobOptions
	pos int
	buf bytes.Buffer
}

func newLobInDescr(tr transform.Transformer, rd io.Reader) *LobInDescr {
	if tr != nil { // cesu8Encoder
		rd = transform.NewReader(rd, tr)
	}
	return &LobInDescr{rd: rd}
}

func (d *LobInDescr) String() string {
	// restrict output size
	return fmt.Sprintf("options %s size %d pos %d bytes %v", d.opt, d.buf.Len(), d.pos, d.buf.Bytes()[:min(d.buf.Len(), 25)])
}

// IsLastData returns true in case of last data package read, false otherwise.
func (d *LobInDescr) IsLastData() bool { return d.opt.isLastData() }

// FetchNext fetches the next lob chunk.
func (d *LobInDescr) FetchNext(chunkSize int) error {
	/*
		We need to guarantee, that a max amount of data is read to prevent
		piece wise LOB writing when avoidable
		--> copy up to chunkSize
	*/
	d.buf.Reset()
	_, err := io.CopyN(&d.buf, d.rd, int64(chunkSize))
	d.opt = loDataincluded
	if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return err
	}
	d.opt |= loLastdata
	return nil
}

func (d *LobInDescr) setPos(pos int) { d.pos = pos }

func (d *LobInDescr) size() int { return d.buf.Len() }

func (d *LobInDescr) writeFirst(enc *encoding.Encoder) { enc.Bytes(d.buf.Bytes()) }

// LocatorID represents a locotor id.
type LocatorID uint64 // byte[locatorIdSize]

// LobReader is the interface for reading lob streams.
type LobReader interface {
	ReadLob(request *ReadLobRequest, reply *ReadLobReply) error
}

var lobOutDescrPool = sync.Pool{New: func() any { return new(lobOutDescr) }}

// lobOutDescr represents a lob output descriptor.
type lobOutDescr struct {
	// if set -> char based
	tr transform.Transformer
	/*
	   readFn is set by decode if additional data packages need to be read (not last data)
	*/
	lobReader LobReader
	chunkSize int
	/*
		HDB does not return lob type code but undefined only
		--> ltc is always ltcUndefined
		--> use isCharBased instead of type code check
	*/
	ltc     lobTypecode
	opt     lobOptions
	numChar int64
	numByte int64
	id      LocatorID
	b       []byte

	// scan attributes.
	wr         io.Writer
	lobRequest *ReadLobRequest
	lobReply   *ReadLobReply
}

func newLobOutDescr(tr transform.Transformer, lobReader LobReader, chunkSize int) *lobOutDescr {
	descr := lobOutDescrPool.Get().(*lobOutDescr)
	descr.tr = tr
	descr.lobReader = lobReader
	descr.chunkSize = chunkSize
	return descr
}

func (d *lobOutDescr) String() string {
	return fmt.Sprintf("typecode %s options %s numChar %d numByte %d id %d bytes %v", d.ltc, d.opt, d.numChar, d.numByte, d.id, d.b)
}

func (d *lobOutDescr) decode(dec *encoding.Decoder) bool {
	d.ltc = lobTypecode(dec.Int8())
	d.opt = lobOptions(dec.Int8())
	if d.opt.isNull() {
		return true
	}
	dec.Skip(2)
	d.numChar = dec.Int64()
	d.numByte = dec.Int64()
	d.id = LocatorID(dec.Uint64())
	size := int(dec.Int32())
	d.b = slices.Grow(d.b, size)[:size]
	dec.Bytes(d.b)
	return false
}

func (d *lobOutDescr) countChars(b []byte) (numChar int) {
	s := unsafe.ByteSlice2String(b)
	for _, r := range s {
		numChar++
		if utf8.RuneLen(r) == 4 {
			numChar++ // caution: hdb counts 2 chars in case of surrogate pair
		}
	}
	return
}

func (d *lobOutDescr) write(b []byte) (int, error) {
	if d.tr == nil {
		if _, err := d.wr.Write(b); err != nil {
			return len(b), err
		}
		return len(b), nil
	}
	d.tr.Reset()
	// cesu8 -> utf8 (always enough space)
	nDst, _, err := d.tr.Transform(b, b, false)
	if err != nil && err != transform.ErrShortSrc { //nolint: errorlint
		return nDst, err
	}

	numChar := d.countChars(b[:nDst])
	if _, err := d.wr.Write(b[:nDst]); err != nil {
		return numChar, err
	}
	return numChar, nil
}

func (d *lobOutDescr) scan(wr io.Writer) error {
	d.wr = wr

	numChar, err := d.write(d.b)
	if err != nil {
		return err
	}

	if d.opt.isLastData() {
		return nil
	}

	if d.lobRequest == nil {
		d.lobRequest = new(ReadLobRequest)
	}
	if d.lobReply == nil {
		d.lobReply = &ReadLobReply{lobOutDescr: d}
	}
	d.lobRequest.id = d.id
	d.lobRequest.ofs = int64(numChar)
	d.lobRequest.chunkSize = d.chunkSize
	return d.lobReader.ReadLob(d.lobRequest, d.lobReply)
}

// Scan implements the LobScanner interface.
func (d *lobOutDescr) Scan(wr io.Writer) error {
	err := d.scan(wr)
	// if the writer is a pipe-end -> close at the end
	if pwr, ok := wr.(*io.PipeWriter); ok {
		if err != nil {
			pwr.CloseWithError(err)
		} else {
			pwr.Close()
		}
	}
	lobOutDescrPool.Put(d)
	return err
}

func (d *lobOutDescr) Write() (int, error) {
	n, err := d.write(d.b)
	if err != nil {
		return n, err
	}
	if d.opt.isLastData() {
		return n, io.EOF
	}
	d.lobRequest.ofs += int64(n)
	return n, nil
}

/*
write lobs:
- write lob field to database in chunks
- loop:
  - writeLobRequest
  - writeLobReply
*/

// WriteLobDescr represents a lob descriptor for writes (lob -> db).
type WriteLobDescr struct {
	LobInDescr *LobInDescr
	ID         LocatorID
	opt        lobOptions
	ofs        int64
	b          []byte
}

func (d WriteLobDescr) String() string {
	return fmt.Sprintf("id %d options %s offset %d bytes %v", d.ID, d.opt, d.ofs, d.b)
}

// IsLastData returns true in case of last data package read, false otherwise.
func (d *WriteLobDescr) IsLastData() bool { return d.opt.isLastData() }

// FetchNext fetches the next lob chunk.
func (d *WriteLobDescr) FetchNext(chunkSize int) error {
	if err := d.LobInDescr.FetchNext(chunkSize); err != nil {
		return err
	}
	d.opt = d.LobInDescr.opt
	d.ofs = -1 // offset (-1 := append)
	d.b = d.LobInDescr.buf.Bytes()
	return nil
}

// sniffer.
func (d *WriteLobDescr) decode(dec *encoding.Decoder) error {
	d.ID = LocatorID(dec.Uint64())
	d.opt = lobOptions(dec.Int8())
	d.ofs = dec.Int64()
	size := dec.Int32()
	d.b = make([]byte, size)
	dec.Bytes(d.b)
	return nil
}

// write chunk to db.
func (d *WriteLobDescr) encode(enc *encoding.Encoder) error {
	enc.Uint64(uint64(d.ID))
	enc.Int8(int8(d.opt))
	enc.Int64(d.ofs)
	enc.Int32(int32(len(d.b))) //nolint: gosec
	enc.Bytes(d.b)
	return nil
}

// WriteLobRequest represents a lob write request part.
type WriteLobRequest struct {
	Descrs []*WriteLobDescr
}

func (r *WriteLobRequest) String() string { return fmt.Sprintf("descriptors %v", r.Descrs) }

func (r *WriteLobRequest) size() int {
	size := 0
	for _, descr := range r.Descrs {
		size += (writeLobRequestSize + len(descr.b))
	}
	return size
}

func (r *WriteLobRequest) numArg() int { return len(r.Descrs) }

// sniffer.
func (r *WriteLobRequest) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	r.Descrs = make([]*WriteLobDescr, numArg)
	for i := range numArg {
		r.Descrs[i] = &WriteLobDescr{}
		if err := r.Descrs[i].decode(dec); err != nil {
			return err
		}
	}
	return nil
}

func (r *WriteLobRequest) encode(enc *encoding.Encoder) error {
	for _, descr := range r.Descrs {
		if err := descr.encode(enc); err != nil {
			return err
		}
	}
	return nil
}

// WriteLobReply represents a lob write reply part.
type WriteLobReply struct {
	// write lob fields to db (reply)
	// - returns ids which have not been written completely
	IDs []LocatorID
}

func (r *WriteLobReply) String() string { return fmt.Sprintf("ids %v", r.IDs) }

func (r *WriteLobReply) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	r.IDs = resizeSlice(r.IDs, numArg)

	for i := range numArg {
		r.IDs[i] = LocatorID(dec.Uint64())
	}
	return dec.Error()
}

// ReadLobRequest represents a lob read request part.
type ReadLobRequest struct {
	/*
	   read lobs:
	   - read lob field from database in chunks
	   - loop:
	     - readLobRequest
	     - readLobReply

	   - read lob reply
	     seems like readLobreply returns only a result for one lob - even if more then one is requested
	     --> read single lobs
	*/
	id        LocatorID
	ofs       int64
	chunkSize int
}

func (r *ReadLobRequest) String() string {
	return fmt.Sprintf("id %d offset %d size %d", r.id, r.ofs, r.chunkSize)
}

// sniffer.
func (r *ReadLobRequest) decode(dec *encoding.Decoder) error {
	r.id = LocatorID(dec.Uint64())
	r.ofs = dec.Int64()
	r.chunkSize = int(dec.Int32())
	dec.Skip(4)
	return nil
}

func (r *ReadLobRequest) encode(enc *encoding.Encoder) error {
	enc.Uint64(uint64(r.id))
	enc.Int64(r.ofs + 1)          // 1-based
	enc.Int32(int32(r.chunkSize)) //nolint: gosec
	enc.Zeroes(4)
	return nil
}

// ReadLobReply represents a lob read reply part.
type ReadLobReply struct {
	*lobOutDescr
}

func (r *ReadLobReply) String() string {
	return fmt.Sprintf("id %d options %s bytes %v", r.id, r.opt, r.b)
}

// needed if instantiated generically (e.g.sniffer).
func (r *ReadLobReply) init() {
	r.lobOutDescr = new(lobOutDescr)
}

func (r *ReadLobReply) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	if numArg != 1 {
		panic("numArg == 1 expected")
	}
	id := LocatorID(dec.Uint64())
	if id != r.id {
		return fmt.Errorf("invalid locator id %d - expected %d", id, r.id)
	}
	r.opt = lobOptions(dec.Int8())
	size := int(dec.Int32())
	dec.Skip(3)
	r.b = slices.Grow(r.b, size)[:size]
	dec.Bytes(r.b)
	return nil
}
