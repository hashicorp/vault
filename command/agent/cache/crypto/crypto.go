package crypto

import (
	"context"
)

const (
	KeyID = "root"
)

// KeyManager TODO
type KeyManager interface {
	Get() []byte
	Renewable() bool
	Renewer(context.Context, chan struct{}) error
	Encrypt(context.Context, []byte, []byte) ([]byte, error)
	Decrypt(context.Context, []byte, []byte) ([]byte, error)
}
