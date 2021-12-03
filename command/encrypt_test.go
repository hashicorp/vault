package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncryptPassphrase(t *testing.T) {
	passphrase := []byte("my excellent password")
	data := []byte("my secret data")

	e := EncryptCommand{}
	encrypted, err := e.encrypt(data, passphrase, true)
	assert.NoError(t, err)

	decrypted, err := e.decrypt(encrypted, passphrase, true)
	assert.NoError(t, err)

	assert.Equal(t, data, decrypted)
}

func Test_EncryptKey(t *testing.T) {
	key := "a key that must be 32 bytes long"
	data := []byte("my secret data")

	e := EncryptCommand{}
	encrypted, err := e.encrypt(data, []byte(key), false)
	assert.NoError(t, err)

	decrypted, err := e.decrypt(encrypted, []byte(key), false)
	assert.NoError(t, err)

	assert.Equal(t, data, decrypted)
}
