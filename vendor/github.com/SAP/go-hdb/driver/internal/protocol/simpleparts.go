package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
	"github.com/SAP/go-hdb/driver/unicode/cesu8"
)

// ClientID represents a client id part.
type ClientID []byte

func (id ClientID) String() string { return string(id) }
func (id ClientID) size() int      { return len(id) }
func (id *ClientID) decodeBufLen(dec *encoding.Decoder, bufLen int) error {
	*id = resizeSlice(*id, bufLen)
	dec.Bytes(*id)
	return dec.Error()
}
func (id ClientID) encode(enc *encoding.Encoder) error { enc.Bytes(id); return nil }

// Command represents a command part with cesu8 content.
type Command []byte

func (c Command) String() string { return string(c) }
func (c Command) size() int      { return cesu8.Size(c) }
func (c *Command) decodeBufLen(dec *encoding.Decoder, bufLen int) error {
	*c = resizeSlice(*c, bufLen)
	var err error
	*c, err = dec.CESU8Bytes(len(*c))
	if err != nil {
		return err
	}
	return dec.Error()
}
func (c Command) encode(enc *encoding.Encoder) error { _, err := enc.CESU8Bytes(c); return err }

// Fetchsize represents a fetch size part.
type Fetchsize int32

func (s Fetchsize) String() string { return fmt.Sprintf("fetchsize %d", s) }
func (s *Fetchsize) decode(dec *encoding.Decoder) error {
	*s = Fetchsize(dec.Int32())
	return dec.Error()
}
func (s Fetchsize) encode(enc *encoding.Encoder) error { enc.Int32(int32(s)); return nil }

// StatementID represents the statement id part type.
type StatementID uint64

func (id StatementID) String() string { return fmt.Sprintf("%d", id) }

// Decode implements the partDecoder interface.
func (id *StatementID) decode(dec *encoding.Decoder) error {
	*id = StatementID(dec.Uint64())
	return dec.Error()
}

// Encode implements the partEncoder interface.
func (id StatementID) encode(enc *encoding.Encoder) error { enc.Uint64(uint64(id)); return nil }
