package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

type clientInfo map[string]string

func (c clientInfo) String() string { return fmt.Sprintf("%v", map[string]string(c)) }

func (c clientInfo) size() int {
	size := 0
	for k, v := range c {
		size += encoding.Cesu8FieldSize(k)
		size += encoding.Cesu8FieldSize(v)
	}
	return size
}

func (c clientInfo) numArg() int { return len(c) }

func (c *clientInfo) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	*c = clientInfo{} // no reuse of maps - create new one

	for range numArg {
		k, err := dec.Cesu8Field()
		if err != nil {
			return err
		}
		v, err := dec.Cesu8Field()
		if err != nil {
			return err
		}
		// set key value
		(*c)[string(k.([]byte))] = string(v.([]byte))
	}
	return dec.Error()
}

func (c clientInfo) encode(enc *encoding.Encoder) error {
	for k, v := range c {
		if err := enc.Cesu8Field(k); err != nil {
			return err
		}
		if err := enc.Cesu8Field(v); err != nil {
			return err
		}
	}
	return nil
}
