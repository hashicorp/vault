package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encrypt(t *testing.T) {
	passphrase := "my excellent password"
	data := []byte("my secret data")

	e := EncryptCommand{}
	encrypted, err := e.encrypt(data, passphrase)
	assert.NoError(t, err)

	d := DecryptCommand{}
	decrypted, err := d.decrypt(encrypted, passphrase)
	assert.NoError(t, err)

	assert.Equal(t, data, decrypted)
}
