package cacheboltdb

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/dhutil"
)

// AES is a simple encryption implementation
type AES struct {
	key    []byte
	logger hclog.Logger
	aad    []byte
}

// AESConfig parameters
type AESConfig struct {
	Key    []byte
	Logger hclog.Logger
	AAD    []byte
}

// NewAES returns a new AES object with some sanity checking
func NewAES(config *AESConfig) (*AES, error) {
	if len(config.Key) != 32 {
		return nil, fmt.Errorf("key length is %d but must be 32", len(config.Key))
	}
	if config.Logger == nil {
		config.Logger = hclog.Default()
	}
	return &AES{
		key:    config.Key,
		logger: config.Logger,
		aad:    config.AAD,
	}, nil
}

// Encrypt accepts plaintext and encrypts into a json-encoded dhutil.Envelope
func (a *AES) Encrypt(plainText []byte) ([]byte, error) {
	var err error
	resp := new(dhutil.Envelope)
	resp.EncryptedPayload, resp.Nonce, err = dhutil.EncryptAES(a.key, plainText, a.aad)
	if err != nil {
		return nil, err
	}
	m, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Decrypt accepts cipherText in the format of dhutil.Envelope and decrypts
func (a *AES) Decrypt(cipherText []byte) ([]byte, error) {
	resp := new(dhutil.Envelope)
	if err := json.Unmarshal(cipherText, resp); err != nil {
		return nil, err
	}
	plainText, err := dhutil.DecryptAES(a.key, resp.EncryptedPayload, resp.Nonce, a.aad)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
