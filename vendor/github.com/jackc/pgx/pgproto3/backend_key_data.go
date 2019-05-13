package pgproto3

import (
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type BackendKeyData struct {
	ProcessID uint32
	SecretKey uint32
}

func (*BackendKeyData) Backend() {}

func (dst *BackendKeyData) Decode(src []byte) error {
	if len(src) != 8 {
		return &invalidMessageLenErr{messageType: "BackendKeyData", expectedLen: 8, actualLen: len(src)}
	}

	dst.ProcessID = binary.BigEndian.Uint32(src[:4])
	dst.SecretKey = binary.BigEndian.Uint32(src[4:])

	return nil
}

func (src *BackendKeyData) Encode(dst []byte) []byte {
	dst = append(dst, 'K')
	dst = pgio.AppendUint32(dst, 12)
	dst = pgio.AppendUint32(dst, src.ProcessID)
	dst = pgio.AppendUint32(dst, src.SecretKey)
	return dst
}

func (src *BackendKeyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string
		ProcessID uint32
		SecretKey uint32
	}{
		Type:      "BackendKeyData",
		ProcessID: src.ProcessID,
		SecretKey: src.SecretKey,
	})
}
