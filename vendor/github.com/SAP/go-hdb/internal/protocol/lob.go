// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"
	"io"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

const (
	writeLobRequestSize = 21
)

// variable (unit testing)
//var lobChunkSize = 1 << 14 //TODO: check size
//var lobChunkSize int32 = 4096 //TODO: check size

// lob options
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

func (o lobOptions) String() string {
	t := make([]string, 0, len(lobOptionsText))

	for option, text := range lobOptionsText {
		if (o & option) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

func (o lobOptions) isLastData() bool { return (o & loLastdata) != 0 }
func (o lobOptions) isNull() bool     { return (o & loNullindicator) != 0 }

//go:generate stringer -type=lobTypecode
// lob typecode
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

// WriterSetter is the interface wrapping the SetWriter method (Lob handling).
type WriterSetter interface{ SetWriter(w io.Writer) error }

// sessionSetter is the interface wrapping the setSession method (lob handling).
type sessionSetter interface{ setSession(s *Session) }

var _ WriterSetter = (*lobOutDescr)(nil)
var _ sessionSetter = (*lobOutDescr)(nil)

/*
TODO description
lobOutDescr

*/
type lobInDescr struct {
	/*
		currently no data is transformed for input parameters
		--> opt == 0 (no data included)
		--> size == 0
		--> pos == 0
		--> b == nil
	*/
	opt  lobOptions
	size int32
	pos  int32
	b    []byte // currently no data is transformed for input parameters
}

/*
TODO description
lobOutDescr

*/
type lobOutDescr struct {
	s           *Session
	isCharBased bool
	/*
		HDB does not return lob type code but undefined only
		--> ltc is always ltcUndefined
		--> use isCharBased instead of type code check
	*/
	ltc     lobTypecode
	opt     lobOptions
	numChar int64
	numByte int64
	id      locatorID
	b       []byte
}

func (d *lobOutDescr) String() string {
	return fmt.Sprintf("typecode %s options %s numChar %d numByte %d id %d bytes %v", d.ltc, d.opt, d.numChar, d.numByte, d.id, d.b)
}
func (d *lobOutDescr) setSession(s *Session) { d.s = s }

// SetWriter implements the WriterSetter interface.
func (d *lobOutDescr) SetWriter(wr io.Writer) error { return d.s.decodeLobs(d, wr) }

/*
write lobs:
- write lob field to database in chunks
- loop:
  - writeLobRequest
  - writeLobReply
*/

// descriptor for writes (lob -> db)
type writeLobDescr struct {
	id  locatorID
	opt lobOptions
	ofs int64
	b   []byte
}

func (d writeLobDescr) String() string {
	return fmt.Sprintf("id %d options %s offset %d bytes %v", d.id, d.opt, d.ofs, d.b)
}

// sniffer
func (d *writeLobDescr) decode(dec *encoding.Decoder) error {
	d.id = locatorID(dec.Uint64())
	d.opt = lobOptions(dec.Int8())
	d.ofs = dec.Int64()
	size := dec.Int32()
	d.b = make([]byte, size)
	dec.Bytes(d.b)
	return nil
}

// write chunk to db
func (d *writeLobDescr) encode(enc *encoding.Encoder) error {
	enc.Uint64(uint64(d.id))
	enc.Int8(int8(d.opt))
	enc.Int64(d.ofs)
	enc.Int32(int32(len(d.b)))
	enc.Bytes(d.b)
	return nil
}

// write lob fields to db (request)
type writeLobRequest struct {
	descrs []*writeLobDescr
}

func (r *writeLobRequest) String() string { return fmt.Sprintf("descriptors %v", r.descrs) }

func (r *writeLobRequest) size() int {
	size := 0
	for _, descr := range r.descrs {
		size += (writeLobRequestSize + len(descr.b))
	}
	return size
}

func (r *writeLobRequest) numArg() int {
	return len(r.descrs)
}

// sniffer
func (r *writeLobRequest) decode(dec *encoding.Decoder, ph *partHeader) error {
	numArg := ph.numArg()
	r.descrs = make([]*writeLobDescr, numArg)
	for i := 0; i < numArg; i++ {
		r.descrs[i] = &writeLobDescr{}
		if err := r.descrs[i].decode(dec); err != nil {
			return err
		}
	}
	return nil
}

func (r *writeLobRequest) encode(enc *encoding.Encoder) error {
	for _, descr := range r.descrs {
		if err := descr.encode(enc); err != nil {
			return err
		}
	}
	return nil
}

// write lob fields to db (reply)
// - returns ids which have not been written completely
type writeLobReply struct {
	ids []locatorID
}

func (r *writeLobReply) String() string { return fmt.Sprintf("ids %v", r.ids) }

func (r *writeLobReply) reset(numArg int) {
	if r.ids == nil || cap(r.ids) < numArg {
		r.ids = make([]locatorID, numArg)
	} else {
		r.ids = r.ids[:numArg]
	}
}

func (r *writeLobReply) decode(dec *encoding.Decoder, ph *partHeader) error {
	numArg := ph.numArg()
	r.reset(numArg)

	for i := 0; i < numArg; i++ {
		r.ids[i] = locatorID(dec.Uint64())
	}
	return dec.Error()
}

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

type readLobRequest struct {
	id        locatorID
	ofs       int64
	chunkSize int32
}

func (r *readLobRequest) String() string {
	return fmt.Sprintf("id %d offset %d size %d", r.id, r.ofs, r.chunkSize)
}

// sniffer
func (r *readLobRequest) decode(dec *encoding.Decoder, ph *partHeader) error {
	r.id = locatorID(dec.Uint64())
	r.ofs = dec.Int64()
	r.chunkSize = dec.Int32()
	dec.Skip(4)
	return nil
}

func (r *readLobRequest) encode(enc *encoding.Encoder) error {
	enc.Uint64(uint64(r.id))
	enc.Int64(r.ofs + 1) //1-based
	enc.Int32(r.chunkSize)
	enc.Zeroes(4)
	return nil
}

type readLobReply struct {
	id  locatorID
	opt lobOptions
	b   []byte
}

func (r *readLobReply) String() string {
	return fmt.Sprintf("id %d options %s bytes %v", r.id, r.opt, r.b)
}

func (r *readLobReply) resize(size int) {
	if r.b == nil || size > cap(r.b) {
		r.b = make([]byte, size)
	} else {
		r.b = r.b[:size]
	}
}

func (r *readLobReply) decode(dec *encoding.Decoder, ph *partHeader) error {
	if ph.numArg() != 1 {
		panic("numArg == 1 expected")
	}
	r.id = locatorID(dec.Uint64())
	r.opt = lobOptions(dec.Int8())
	size := int(dec.Int32())
	dec.Skip(3)
	r.resize(size)
	dec.Bytes(r.b)
	return nil
}
