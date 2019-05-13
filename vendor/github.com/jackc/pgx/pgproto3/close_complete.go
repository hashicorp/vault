package pgproto3

import (
	"encoding/json"
)

type CloseComplete struct{}

func (*CloseComplete) Backend() {}

func (dst *CloseComplete) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "CloseComplete", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *CloseComplete) Encode(dst []byte) []byte {
	return append(dst, '3', 0, 0, 0, 4)
}

func (src *CloseComplete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "CloseComplete",
	})
}
