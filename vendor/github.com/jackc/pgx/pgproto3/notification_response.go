package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/jackc/pgx/pgio"
)

type NotificationResponse struct {
	PID     uint32
	Channel string
	Payload string
}

func (*NotificationResponse) Backend() {}

func (dst *NotificationResponse) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	pid := binary.BigEndian.Uint32(buf.Next(4))

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	channel := string(b[:len(b)-1])

	b, err = buf.ReadBytes(0)
	if err != nil {
		return err
	}
	payload := string(b[:len(b)-1])

	*dst = NotificationResponse{PID: pid, Channel: channel, Payload: payload}
	return nil
}

func (src *NotificationResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'A')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Channel...)
	dst = append(dst, 0)
	dst = append(dst, src.Payload...)
	dst = append(dst, 0)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *NotificationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string
		PID     uint32
		Channel string
		Payload string
	}{
		Type:    "NotificationResponse",
		PID:     src.PID,
		Channel: src.Channel,
		Payload: src.Payload,
	})
}
