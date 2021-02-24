package crypto

import (
	"context"
)

const (
	KeyID = "root"
)

type KeyManager interface {
	GetKey() []byte
	GetPersistentKey() ([]byte, error)
	Renewable() bool
	Renewer(context.Context) error
	Encrypt(context.Context, []byte, []byte) ([]byte, error)
	Decrypt(context.Context, []byte, []byte) ([]byte, error)
}
