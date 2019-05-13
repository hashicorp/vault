package pgproto3

import (
	"encoding/json"
)

type ReadyForQuery struct {
	TxStatus byte
}

func (*ReadyForQuery) Backend() {}

func (dst *ReadyForQuery) Decode(src []byte) error {
	if len(src) != 1 {
		return &invalidMessageLenErr{messageType: "ReadyForQuery", expectedLen: 1, actualLen: len(src)}
	}

	dst.TxStatus = src[0]

	return nil
}

func (src *ReadyForQuery) Encode(dst []byte) []byte {
	return append(dst, 'Z', 0, 0, 0, 5, src.TxStatus)
}

func (src *ReadyForQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		TxStatus string
	}{
		Type:     "ReadyForQuery",
		TxStatus: string(src.TxStatus),
	})
}
