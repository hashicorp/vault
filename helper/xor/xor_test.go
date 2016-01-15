package xor

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
)

const (
	tokenB64    = "ZGE0N2JiODkzYjhkMDYxYw=="
	xorB64      = "iGiQYG9L0nIp+jRL5+Zk2w=="
	expectedB64 = "7AmkVw0p6ksamAwv19BVuA=="
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

func TestBase64XOR(t *testing.T) {
	ret, err := XORBase64(tokenB64, xorB64)
	if err != nil {
		t.Fatal(err)
	}
	if res := base64.StdEncoding.EncodeToString(ret); res != expectedB64 {
		t.Fatalf("bad: %s", res)
	}
}
