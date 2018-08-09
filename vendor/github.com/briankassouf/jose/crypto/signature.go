package crypto

import (
	"encoding/json"

	"github.com/briankassouf/jose"
)

// Signature is a JWS signature.
type Signature []byte

// MarshalJSON implements json.Marshaler for a signature.
func (s Signature) MarshalJSON() ([]byte, error) {
	return jose.EncodeEscape(s), nil
}

// Base64 helps implements jose.Encoder for Signature.
func (s Signature) Base64() ([]byte, error) {
	return jose.Base64Encode(s), nil
}

// UnmarshalJSON implements json.Unmarshaler for signature.
func (s *Signature) UnmarshalJSON(b []byte) error {
	dec, err := jose.DecodeEscaped(b)
	if err != nil {
		return err
	}
	*s = Signature(dec)
	return nil
}

var (
	_ json.Marshaler   = (Signature)(nil)
	_ json.Unmarshaler = (*Signature)(nil)
	_ jose.Encoder     = (Signature)(nil)
)
