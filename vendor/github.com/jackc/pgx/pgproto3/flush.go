package pgproto3

import (
	"encoding/json"
)

type Flush struct{}

func (*Flush) Frontend() {}

func (dst *Flush) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "Flush", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *Flush) Encode(dst []byte) []byte {
	return append(dst, 'H', 0, 0, 0, 4)
}

func (src *Flush) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "Flush",
	})
}
