package pgproto3

import (
	"bytes"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type Query struct {
	String string
}

func (*Query) Frontend() {}

func (dst *Query) Decode(src []byte) error {
	i := bytes.IndexByte(src, 0)
	if i != len(src)-1 {
		return &invalidMessageFormatErr{messageType: "Query"}
	}

	dst.String = string(src[:i])

	return nil
}

func (src *Query) Encode(dst []byte) []byte {
	dst = append(dst, 'Q')
	dst = pgio.AppendInt32(dst, int32(4+len(src.String)+1))

	dst = append(dst, src.String...)
	dst = append(dst, 0)

	return dst
}

func (src *Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string
		String string
	}{
		Type:   "Query",
		String: src.String,
	})
}
