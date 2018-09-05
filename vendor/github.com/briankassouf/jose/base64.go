package jose

import "encoding/base64"

// Encoder is satisfied if the type can marshal itself into a valid
// structure for a JWS.
type Encoder interface {
	// Base64 implies T -> JSON -> RawURLEncodingBase64
	Base64() ([]byte, error)
}

// Base64Decode decodes a base64-encoded byte slice.
func Base64Decode(b []byte) ([]byte, error) {
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(buf, b)
	return buf[:n], err
}

// Base64Encode encodes a byte slice.
func Base64Encode(b []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(buf, b)
	return buf
}

// EncodeEscape base64-encodes a byte slice but escapes it for JSON.
// It'll return the format: `"base64"`
func EncodeEscape(b []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(b))+2)
	buf[0] = '"'
	base64.RawURLEncoding.Encode(buf[1:], b)
	buf[len(buf)-1] = '"'
	return buf
}

// DecodeEscaped decodes a base64-encoded byte slice straight from a JSON
// structure. It assumes it's in the format: `"base64"`, but can handle
// cases where it's not.
func DecodeEscaped(b []byte) ([]byte, error) {
	if len(b) > 1 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	return Base64Decode(b)
}
