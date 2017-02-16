package jws

import (
	"encoding/json"

	"github.com/SermoDigital/jose"
)

// payload represents the payload of a JWS.
type payload struct {
	v interface{}
	u json.Unmarshaler
	_ struct{}
}

// MarshalJSON implements json.Marshaler for payload.
func (p *payload) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(p.v)
	if err != nil {
		return nil, err
	}
	return jose.EncodeEscape(b), nil
}

// Base64 implements jose.Encoder.
func (p *payload) Base64() ([]byte, error) {
	b, err := json.Marshal(p.v)
	if err != nil {
		return nil, err
	}
	return jose.Base64Encode(b), nil
}

// MarshalJSON implements json.Unmarshaler for payload.
func (p *payload) UnmarshalJSON(b []byte) error {
	b2, err := jose.DecodeEscaped(b)
	if err != nil {
		return err
	}
	if p.u != nil {
		err := p.u.UnmarshalJSON(b2)
		p.v = p.u
		return err
	}
	return json.Unmarshal(b2, &p.v)
}

var (
	_ json.Marshaler   = (*payload)(nil)
	_ json.Unmarshaler = (*payload)(nil)
	_ jose.Encoder     = (*payload)(nil)
)
