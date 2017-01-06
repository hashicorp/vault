package jws

import (
	"sync"

	"github.com/SermoDigital/jose/crypto"
)

var (
	mu sync.RWMutex

	signingMethods = map[string]crypto.SigningMethod{
		crypto.SigningMethodES256.Alg(): crypto.SigningMethodES256,
		crypto.SigningMethodES384.Alg(): crypto.SigningMethodES384,
		crypto.SigningMethodES512.Alg(): crypto.SigningMethodES512,

		crypto.SigningMethodPS256.Alg(): crypto.SigningMethodPS256,
		crypto.SigningMethodPS384.Alg(): crypto.SigningMethodPS384,
		crypto.SigningMethodPS512.Alg(): crypto.SigningMethodPS512,

		crypto.SigningMethodRS256.Alg(): crypto.SigningMethodRS256,
		crypto.SigningMethodRS384.Alg(): crypto.SigningMethodRS384,
		crypto.SigningMethodRS512.Alg(): crypto.SigningMethodRS512,

		crypto.SigningMethodHS256.Alg(): crypto.SigningMethodHS256,
		crypto.SigningMethodHS384.Alg(): crypto.SigningMethodHS384,
		crypto.SigningMethodHS512.Alg(): crypto.SigningMethodHS512,

		crypto.Unsecured.Alg(): crypto.Unsecured,
	}
)

// RegisterSigningMethod registers the crypto.SigningMethod in the global map.
// This is typically done inside the caller's init function.
func RegisterSigningMethod(sm crypto.SigningMethod) {
	alg := sm.Alg()
	if GetSigningMethod(alg) != nil {
		panic("jose/jws: cannot duplicate signing methods")
	}

	if !sm.Hasher().Available() {
		panic("jose/jws: specific hash is unavailable")
	}

	mu.Lock()
	signingMethods[alg] = sm
	mu.Unlock()
}

// RemoveSigningMethod removes the crypto.SigningMethod from the global map.
func RemoveSigningMethod(sm crypto.SigningMethod) {
	mu.Lock()
	delete(signingMethods, sm.Alg())
	mu.Unlock()
}

// GetSigningMethod retrieves a crypto.SigningMethod from the global map.
func GetSigningMethod(alg string) (method crypto.SigningMethod) {
	mu.RLock()
	method = signingMethods[alg]
	mu.RUnlock()
	return method
}
