package crypto

import (
	"crypto"
	"encoding/json"
	"hash"
	"io"
)

func init() {
	crypto.RegisterHash(crypto.Hash(0), h)
}

// h is passed to crypto.RegisterHash.
func h() hash.Hash {
	return &f{Writer: nil}
}

type f struct{ io.Writer }

// Sum helps implement the hash.Hash interface.
func (_ *f) Sum(b []byte) []byte { return nil }

// Reset helps implement the hash.Hash interface.
func (_ *f) Reset() {}

// Size helps implement the hash.Hash interface.
func (_ *f) Size() int { return -1 }

// BlockSize helps implement the hash.Hash interface.
func (_ *f) BlockSize() int { return -1 }

// Unsecured is the default "none" algorithm.
var Unsecured = &SigningMethodNone{
	Name: "none",
	Hash: crypto.Hash(0),
}

// SigningMethodNone is the default "none" algorithm.
type SigningMethodNone struct {
	Name string
	Hash crypto.Hash
	_    struct{}
}

// Verify helps implement the SigningMethod interface.
func (_ *SigningMethodNone) Verify(_ []byte, _ Signature, _ interface{}) error {
	return nil
}

// Sign helps implement the SigningMethod interface.
func (_ *SigningMethodNone) Sign(_ []byte, _ interface{}) (Signature, error) {
	return nil, nil
}

// Alg helps implement the SigningMethod interface.
func (m *SigningMethodNone) Alg() string {
	return m.Name
}

// Hasher helps implement the SigningMethod interface.
func (m *SigningMethodNone) Hasher() crypto.Hash {
	return m.Hash
}

// MarshalJSON implements json.Marshaler.
// See SigningMethodECDSA.MarshalJSON() for information.
func (m *SigningMethodNone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodNone)(nil)
