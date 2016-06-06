package azureservicebus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ComputeHmac256 signs a string message with the given key
func ComputeHmac256(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
