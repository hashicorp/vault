package pgproto3

import (
	"encoding/hex"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type CopyData struct {
	Data []byte
}

func (*CopyData) Backend()  {}
func (*CopyData) Frontend() {}

func (dst *CopyData) Decode(src []byte) error {
	dst.Data = src
	return nil
}

func (src *CopyData) Encode(dst []byte) []byte {
	dst = append(dst, 'd')
	dst = pgio.AppendInt32(dst, int32(4+len(src.Data)))
	dst = append(dst, src.Data...)
	return dst
}

func (src *CopyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Data string
	}{
		Type: "CopyData",
		Data: hex.EncodeToString(src.Data),
	})
}
