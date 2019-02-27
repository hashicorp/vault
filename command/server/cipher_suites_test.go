package server

import (
	"crypto/tls"
	"testing"
)

func TestIsBadCipher(t *testing.T) {
	badCipher := tls.TLS_RSA_WITH_AES_128_CBC_SHA
	if !isBadCipher(badCipher) {
		t.Fatalf("TLS_RSA_WITH_AES_128_CBC_SHA is a bad cipher but has not been detected")
	}
	goodCipher := tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
	if isBadCipher(goodCipher) {
		t.Fatalf("TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 is a good cipher but has been detected")
	}
}
