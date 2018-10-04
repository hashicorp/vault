package pgproto3

import (
	"encoding/json"
)

type BindComplete struct{}

func (*BindComplete) Backend() {}

func (dst *BindComplete) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "BindComplete", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *BindComplete) Encode(dst []byte) []byte {
	return append(dst, '2', 0, 0, 0, 4)
}

func (src *BindComplete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "BindComplete",
	})
}
