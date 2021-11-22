package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64OTPGeneration(t *testing.T) {
	token, errCode, err := GenerateOTP(0)
	assert.Len(t, token, 24)
	assert.Equal(t, errCode, 0)
	assert.Nil(t, err)
}

func TestBase62OTPGeneration(t *testing.T) {
	token, errCode, err := GenerateOTP(20)
	assert.Len(t, token, 20)
	assert.Equal(t, errCode, 0)
	assert.Nil(t, err)
}
