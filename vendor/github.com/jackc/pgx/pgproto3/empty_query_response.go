package pgproto3

import (
	"encoding/json"
)

type EmptyQueryResponse struct{}

func (*EmptyQueryResponse) Backend() {}

func (dst *EmptyQueryResponse) Decode(src []byte) error {
	if len(src) != 0 {
		return &invalidMessageLenErr{messageType: "EmptyQueryResponse", expectedLen: 0, actualLen: len(src)}
	}

	return nil
}

func (src *EmptyQueryResponse) Encode(dst []byte) []byte {
	return append(dst, 'I', 0, 0, 0, 4)
}

func (src *EmptyQueryResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
	}{
		Type: "EmptyQueryResponse",
	})
}
