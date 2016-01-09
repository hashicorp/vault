package xor

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandBytes(length int) ([]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("length must be >= 0")
	}

	buf := make([]byte, length)
	if length == 0 {
		return buf, nil
	}

	n, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, fmt.Errorf("unable to read %d bytes; only read %d", length, n)
	}

	return buf, nil
}

func XORBuffers(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of buffers is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

func XORBase64(a, b string) ([]byte, error) {
	aBytes, err := base64.StdEncoding.DecodeString(a)
	if err != nil {
		return nil, fmt.Errorf("error decoding first base64 value: %v", err)
	}
	if aBytes == nil || len(aBytes) == 0 {
		return nil, fmt.Errorf("decoded first base64 value is nil or empty")
	}

	bBytes, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return nil, fmt.Errorf("error decoding second base64 value: %v", err)
	}
	if bBytes == nil || len(bBytes) == 0 {
		return nil, fmt.Errorf("decoded second base64 value is nil or empty")
	}

	if len(aBytes) != len(bBytes) {
		return nil, fmt.Errorf("decoded values are not same length: %d != %d", len(aBytes), len(bBytes))
	}

	buf := make([]byte, len(aBytes))
	for i, _ := range aBytes {
		buf[i] = aBytes[i] ^ bBytes[i]
	}

	return buf, nil
}
