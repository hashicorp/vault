package root

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/pgpkeys"
	"github.com/hashicorp/vault/sdk/helper/xor"
)

var (
	// ErrNoTokenProvided is returned when an empty string is provided as an error
	ErrNoTokenProvided = errors.New("no token provided")
	// ErrNoOTPorPGPKeyProvided is returned when a non-empty token is provided but not OTP or PGP key
	ErrNoOTPorPGPKeyProvided = errors.New("no otp or pgp key provided")
)

// EncodeToken gets a token and an OTP or a PGP Key and encodes it. If no OTP or PGP key is
// provided, an error will be returned. If an OTP is provided, it must have the same length as the token
func EncodeToken(token, otp, pgpKey string, cleanupFunc func()) (string, error) {
	if len(token) == 0 {
		cleanupFunc()
		return "", ErrNoTokenProvided
	}
	// Get the encoded value first so that if there is an error we don't create
	// the root token.
	if len(otp) > 0 {
		// This function performs decoding checks so rather than decode the OTP,
		// just encode the value we're passing in.
		tokenBytes, err := xor.XORBytes([]byte(otp), []byte(token))
		if err != nil {
			cleanupFunc()
			return "", fmt.Errorf("xor of root token failed: %w", err)
		}
		return base64.RawStdEncoding.EncodeToString(tokenBytes), nil
	} else if len(pgpKey) > 0 {
		_, tokenBytesArr, err := pgpkeys.EncryptShares([][]byte{[]byte(token)}, []string{pgpKey})
		if err != nil {
			cleanupFunc()
			return "", fmt.Errorf("error encrypting new root token: %w", err)
		}
		return base64.StdEncoding.EncodeToString(tokenBytesArr[0]), nil
	} else {
		cleanupFunc()
		return "", ErrNoOTPorPGPKeyProvided
	}
}
