package pgproto3

import (
	"encoding/binary"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

const (
	AuthTypeOk                = 0
	AuthTypeCleartextPassword = 3
	AuthTypeMD5Password       = 5
)

type Authentication struct {
	Type uint32

	// MD5Password fields
	Salt [4]byte
}

func (*Authentication) Backend() {}

func (dst *Authentication) Decode(src []byte) error {
	*dst = Authentication{Type: binary.BigEndian.Uint32(src[:4])}

	switch dst.Type {
	case AuthTypeOk:
	case AuthTypeCleartextPassword:
	case AuthTypeMD5Password:
		copy(dst.Salt[:], src[4:8])
	default:
		return errors.Errorf("unknown authentication type: %d", dst.Type)
	}

	return nil
}

func (src *Authentication) Encode(dst []byte) []byte {
	dst = append(dst, 'R')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)
	dst = pgio.AppendUint32(dst, src.Type)

	switch src.Type {
	case AuthTypeMD5Password:
		dst = append(dst, src.Salt[:]...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
