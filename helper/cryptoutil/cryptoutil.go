package cryptoutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

func Blake2b256Hash(key string) []byte {
	hf, _ := blake2b.New256(nil)

	hf.Write([]byte(key))

	return hf.Sum(nil)
}

// HMACSHA256Hash returns a hex-encoded HMAC-SHA256 value based
// on the provided key and raw value to hash.
func HMACSHA256Hash(key, value string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("invalid HMAC key")
	}
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write([]byte(value))
	return hex.EncodeToString(hm.Sum(nil)), nil
}
