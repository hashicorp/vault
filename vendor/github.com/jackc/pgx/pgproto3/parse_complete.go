package pgproto3

import (
	"encoding/json"
)

type ParseComplete struct{}

func (*ParseComplete) Backend() {}

func (dst *ParseComplete) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "ParseComplete", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *ParseComplete) Encode(dst []byte) []byte {
	return append(dst, '1', 0, 0, 0, 4)
}

func (src *ParseComplete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "ParseComplete",
	})
}
