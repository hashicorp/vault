package cachepersist

import "github.com/hashicorp/go-hclog"

// Encryption is an interface for encrypting and decrypting items in storage
type Encryption interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// PassThru is an encryption implementation that doesn't do anything
type PassThru struct {
	Key    string
	Logger hclog.Logger
}

func NewPassThru(key string, logger hclog.Logger) *PassThru {
	return &PassThru{Key: key, Logger: logger}
}

func (p *PassThru) Encrypt(plainText []byte) ([]byte, error) {
	cipherText := plainText
	p.Logger.Trace("encrypted something!")
	return cipherText, nil
}

func (p *PassThru) Decrypt(cipherText []byte) ([]byte, error) {
	plainText := cipherText
	p.Logger.Trace("decrypted something!")
	return plainText, nil
}
