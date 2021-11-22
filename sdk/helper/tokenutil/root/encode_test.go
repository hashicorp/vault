package root

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/pgpkeys"
	"github.com/stretchr/testify/assert"
)

func TestTokenEncodingDecodingWithOTP(t *testing.T) {
	otpTestCases := []struct {
		token               string
		name                string
		otpLength           int
		expectedEncodingErr string
		expectedDecodingErr string
	}{
		{
			token:               "someToken",
			name:                "test token encoding with base64",
			otpLength:           0,
			expectedEncodingErr: "xor of root token failed: length of byte slices is not equivalent: 24 != 9",
			expectedDecodingErr: "",
		},
		{
			token:               "someToken",
			name:                "test token encoding with base62",
			otpLength:           len("someToken"),
			expectedEncodingErr: "",
			expectedDecodingErr: "",
		},
		{
			token:               "someToken",
			name:                "test token encoding with base62 - wrong otp length",
			otpLength:           len("someToken") + 1,
			expectedEncodingErr: "xor of root token failed: length of byte slices is not equivalent: 10 != 9",
			expectedDecodingErr: "",
		},
		{
			token:               "",
			name:                "test no token to encode",
			otpLength:           0,
			expectedEncodingErr: "no token provided",
			expectedDecodingErr: "",
		},
	}
	for _, otpTestCase := range otpTestCases {
		t.Run(otpTestCase.name, func(t *testing.T) {
			cleanupOnErr := false
			otp, _, err := GenerateOTP(otpTestCase.otpLength)
			if err != nil {
				t.Fatal(err.Error())
			}
			encodedToken, err := EncodeToken(otpTestCase.token, otp, "", func() {
				cleanupOnErr = true
			})
			if err != nil || otpTestCase.expectedDecodingErr != "" {
				assert.EqualError(t, err, otpTestCase.expectedEncodingErr)
				assert.True(t, cleanupOnErr)
				return
			}
			assert.NotEqual(t, otp, encodedToken)
			assert.NotEqual(t, encodedToken, otpTestCase.token)
			decodedToken, err := DecodeToken(encodedToken, otp, len(otp))
			if err != nil || otpTestCase.expectedDecodingErr != "" {
				assert.EqualError(t, err, otpTestCase.expectedDecodingErr)
				assert.True(t, cleanupOnErr)
				return
			}
			assert.Equal(t, otpTestCase.token, decodedToken)
			assert.False(t, cleanupOnErr)
		})
	}
}

func TestTokenEncodingDecodingWithNoOTPorPGPKey(t *testing.T) {
	_, err := EncodeToken("", "", "", func() {})
	assert.ErrorIs(t, err, ErrNoTokenProvided)
}

func TestTokenEncodingWithPGPKey(t *testing.T) {
	token := "someToken"
	encodedToken, err := EncodeToken(token, "", pgpkeys.TestPubKey1, func() {})
	assert.Nil(t, err)
	assert.NotEqual(t, encodedToken, token)
	bb, err := pgpkeys.DecryptBytes(encodedToken, pgpkeys.TestPrivKey1)
	assert.Nil(t, err)
	plaintext := bb.String()
	assert.Equal(t, token, plaintext)
}
