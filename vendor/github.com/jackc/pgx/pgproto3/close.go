package pgproto3

import (
	"bytes"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type Close struct {
	ObjectType byte // 'S' = prepared statement, 'P' = portal
	Name       string
}

func (*Close) Frontend() {}

func (dst *Close) Decode(src []byte) error {
	if len(src) < 2 {
		return &invalidMessageFormatErr{messageType: "Close"}
	}

	dst.ObjectType = src[0]
	rp := 1

	idx := bytes.IndexByte(src[rp:], 0)
	if idx != len(src[rp:])-1 {
		return &invalidMessageFormatErr{messageType: "Close"}
	}

	dst.Name = string(src[rp : len(src)-1])

	return nil
}

func (src *Close) Encode(dst []byte) []byte {
	dst = append(dst, 'C')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.ObjectType)
	dst = append(dst, src.Name...)
	dst = append(dst, 0)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *Close) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string
		ObjectType string
		Name       string
	}{
		Type:       "Close",
		ObjectType: string(src.ObjectType),
		Name:       src.Name,
	})
}
