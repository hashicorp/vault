package xor

import (
	"encoding/base64"
	"testing"
)

const (
	tokenB64    = "ZGE0N2JiODkzYjhkMDYxYw=="
	xorB64      = "iGiQYG9L0nIp+jRL5+Zk2w=="
	expectedB64 = "7AmkVw0p6ksamAwv19BVuA=="
)

func TestBase64XOR(t *testing.T) {
	ret, err := XORBase64(tokenB64, xorB64)
	if err != nil {
		t.Fatal(err)
	}
	if res := base64.StdEncoding.EncodeToString(ret); res != expectedB64 {
		t.Fatalf("bad: %s", res)
	}
}
