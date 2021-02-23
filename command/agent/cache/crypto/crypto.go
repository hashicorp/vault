package crypto

import (
	"context"
)

const (
	KeyID = "root"
)

type KeyManager interface {
	GetKey() []byte
	GetPersistentKey() []byte
	Renewable() bool
	Renewer(context.Context, chan struct{}) error
	Encrypt(context.Context, []byte, []byte) ([]byte, error)
	Decrypt(context.Context, []byte, []byte) ([]byte, error)
}
