package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encrypt(t *testing.T) {
	passphrase := "my excellent password"
	data := []byte("my secret data")

	e := EncryptCommand{}
	encrypted, err := e.encrypt(data, passphrase, nil)
	assert.NoError(t, err)

	decrypted, err := e.decrypt(encrypted, passphrase, nil)
	assert.NoError(t, err)

	assert.Equal(t, data, decrypted)
}
