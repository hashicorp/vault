package rand

import (
	"encoding/hex"
	"io"
)

const dash byte = '-'

// UUIDIdempotencyToken provides a utility to get idempotency tokens in the
// UUID format.
type UUIDIdempotencyToken struct {
	uuid *UUID
}

// NewUUIDIdempotencyToken returns a idempotency token provider returning
// tokens in the UUID random format using the reader provided.
func NewUUIDIdempotencyToken(r io.Reader) *UUIDIdempotencyToken {
	return &UUIDIdempotencyToken{uuid: NewUUID(r)}
}

// GetIdempotencyToken returns a random UUID value for Idempotency token.
func (u UUIDIdempotencyToken) GetIdempotencyToken() (string, error) {
	return u.uuid.GetUUID()
}

// UUID provides computing random UUID version 4 values from a random source
// reader.
type UUID struct {
	randSrc io.Reader
}

// NewUUID returns an initialized UUID value that can be used to retrieve
// random UUID values.
func NewUUID(r io.Reader) *UUID {
	return &UUID{randSrc: r}
}

// GetUUID returns a UUID  random string sourced from the random reader the
// UUID was created with. Returns an error if unable to compute the UUID.
func (r *UUID) GetUUID() (string, error) {
	var b [16]byte
	if _, err := io.ReadFull(r.randSrc, b[:]); err != nil {
		return "", err
	}

	return uuidVersion4(b), nil
}

// uuidVersion4 returns a random UUID version 4 from the byte slice provided.
func uuidVersion4(u [16]byte) string {
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29

	// 13th character is "4"
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	// 17th character is "8", "9", "a", or "b"
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10

	var scratch [36]byte

	hex.Encode(scratch[:8], u[0:4])
	scratch[8] = dash
	hex.Encode(scratch[9:13], u[4:6])
	scratch[13] = dash
	hex.Encode(scratch[14:18], u[6:8])
	scratch[18] = dash
	hex.Encode(scratch[19:23], u[8:10])
	scratch[23] = dash
	hex.Encode(scratch[24:], u[10:])

	return string(scratch[:])
}
