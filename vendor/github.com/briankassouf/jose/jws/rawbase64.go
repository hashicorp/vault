package jws

import "encoding/json"

type rawBase64 []byte

// MarshalJSON implements json.Marshaler for rawBase64.
func (r rawBase64) MarshalJSON() ([]byte, error) {
	buf := make([]byte, len(r)+2)
	buf[0] = '"'
	copy(buf[1:], r)
	buf[len(buf)-1] = '"'
	return buf, nil
}

// MarshalJSON implements json.Unmarshaler for rawBase64.
func (r *rawBase64) UnmarshalJSON(b []byte) error {
	if len(b) > 1 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	*r = rawBase64(b)
	return nil
}

var (
	_ json.Marshaler   = (rawBase64)(nil)
	_ json.Unmarshaler = (*rawBase64)(nil)
)
