package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type Execute struct {
	Portal  string
	MaxRows uint32
}

func (*Execute) Frontend() {}

func (dst *Execute) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Portal = string(b[:len(b)-1])

	if buf.Len() < 4 {
		return &invalidMessageFormatErr{messageType: "Execute"}
	}
	dst.MaxRows = binary.BigEndian.Uint32(buf.Next(4))

	return nil
}

func (src *Execute) Encode(dst []byte) []byte {
	dst = append(dst, 'E')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Portal...)
	dst = append(dst, 0)

	dst = pgio.AppendUint32(dst, src.MaxRows)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *Execute) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string
		Portal  string
		MaxRows uint32
	}{
		Type:    "Execute",
		Portal:  src.Portal,
		MaxRows: src.MaxRows,
	})
}
