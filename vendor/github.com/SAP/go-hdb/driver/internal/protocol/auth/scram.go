package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const (
	clientChallengeSize = 64
	serverChallengeSize = 48
	saltSize            = 16
	clientProofSize     = 32
)

func checkSalt(salt []byte) error {
	if len(salt) != saltSize {
		return fmt.Errorf("invalid salt size %d - expected %d", len(salt), saltSize)
	}
	return nil
}

func checkServerChallenge(serverChallenge []byte) error {
	if len(serverChallenge) != serverChallengeSize {
		return fmt.Errorf("invalid server challenge size %d - expected %d", len(serverChallenge), serverChallengeSize)
	}
	return nil
}

func clientChallenge() []byte {
	r := make([]byte, clientChallengeSize)
	if _, err := rand.Read(r); err != nil {
		panic(err)
	}
	return r
}

func clientProof(key, salt, serverChallenge, clientChallenge []byte) ([]byte, error) {
	if len(key) != clientProofSize {
		return nil, fmt.Errorf("invalid key size %d - expected %d", len(key), clientProofSize)
	}
	sig := _hmac(_sha256(key), salt, serverChallenge, clientChallenge)
	if len(sig) != clientProofSize {
		return nil, fmt.Errorf("invalid sig size %d - expected %d", len(key), clientProofSize)
	}
	// xor sig and key into sig (inline: no further allocation).
	for i, v := range key {
		sig[i] ^= v
	}
	return sig, nil
}

func _sha256(p []byte) []byte {
	hash := sha256.New()
	hash.Write(p)
	return hash.Sum(nil)
}

func _hmac(key []byte, prms ...[]byte) []byte {
	hash := hmac.New(sha256.New, key)
	for _, p := range prms {
		hash.Write(p)
	}
	return hash.Sum(nil)
}
