package roottoken

import (
	"testing"

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
			otp, err := GenerateOTP(otpTestCase.otpLength)
			if err != nil {
				t.Fatal(err.Error())
			}
			encodedToken, err := EncodeToken(otpTestCase.token, otp)
			if err != nil || otpTestCase.expectedDecodingErr != "" {
				assert.EqualError(t, err, otpTestCase.expectedEncodingErr)
				return
			}
			assert.NotEqual(t, otp, encodedToken)
			assert.NotEqual(t, encodedToken, otpTestCase.token)
			decodedToken, err := DecodeToken(encodedToken, otp, len(otp))
			if err != nil || otpTestCase.expectedDecodingErr != "" {
				assert.EqualError(t, err, otpTestCase.expectedDecodingErr)
				return
			}
			assert.Equal(t, otpTestCase.token, decodedToken)
		})
	}
}

func TestTokenEncodingDecodingWithNoOTPorPGPKey(t *testing.T) {
	_, err := EncodeToken("", "")
	assert.EqualError(t, err, "no token provided")
}
