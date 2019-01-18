package pgproto3

import (
	"bytes"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type PasswordMessage struct {
	Password string
}

func (*PasswordMessage) Frontend() {}

func (dst *PasswordMessage) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Password = string(b[:len(b)-1])

	return nil
}

func (src *PasswordMessage) Encode(dst []byte) []byte {
	dst = append(dst, 'p')
	dst = pgio.AppendInt32(dst, int32(4+len(src.Password)+1))

	dst = append(dst, src.Password...)
	dst = append(dst, 0)

	return dst
}

func (src *PasswordMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Password string
	}{
		Type:     "PasswordMessage",
		Password: src.Password,
	})
}
