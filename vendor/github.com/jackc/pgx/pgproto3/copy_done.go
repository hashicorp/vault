package pgproto3

import (
	"encoding/json"
)

type CopyDone struct {
}

func (*CopyDone) Backend() {}

func (dst *CopyDone) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "CopyDone", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *CopyDone) Encode(dst []byte) []byte {
	return append(dst, 'c', 0, 0, 0, 4)
}

func (src *CopyDone) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "CopyDone",
	})
}
