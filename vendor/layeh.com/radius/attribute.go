package radius

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"time"
)

// ErrNoAttribute is returned when an attribute was not found when one was
// expected.
var ErrNoAttribute = errors.New("radius: attribute not found")

// Attribute is a RADIUS attribute value, as encoded on the wire.
type Attribute []byte

// Integer returns the given attribute as an integer. An error is returned  if
// the attribute is not 4 bytes long.
func Integer(a Attribute) (uint32, error) {
	if len(a) != 4 {
		return 0, errors.New("invalid length")
	}
	return binary.BigEndian.Uint32(a), nil
}

// NewInteger creates a new Attribute from the given integer value.
func NewInteger(i uint32) Attribute {
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(v, i)
	return Attribute(v)
}

// String returns the given attribute as a string.
func String(a Attribute) string {
	return string(a)
}

// NewString returns a new Attribute from the given string. An error is returned
// if the string length is greater than 253.
func NewString(s string) (Attribute, error) {
	if len(s) > 253 {
		return nil, errors.New("string too long")
	}
	return Attribute(s), nil
}

// Bytes returns the given Attribute as a byte slice.
func Bytes(a Attribute) []byte {
	b := make([]byte, len(a))
	copy(b, []byte(a))
	return b
}

// NewBytes returns a new Attribute from the given byte slice. An error is
// returned if the slice is longer than 253.
func NewBytes(b []byte) (Attribute, error) {
	if len(b) > 253 {
		return nil, errors.New("value too long")
	}
	a := make(Attribute, len(b))
	copy(a, Attribute(b))
	return a, nil
}

// IPAddr returns the given Attribute as an IPv4 IP address. An error is
// returned if the attribute is not 4 bytes long.
func IPAddr(a Attribute) (net.IP, error) {
	if len(a) != net.IPv4len {
		return nil, errors.New("invalid length")
	}
	b := make([]byte, net.IPv4len)
	copy(b, []byte(a))
	return b, nil
}

// NewIPAddr returns a new Attribute from the given IP address. An error is
// returned if the given address is not an IPv4 address.
func NewIPAddr(a net.IP) (Attribute, error) {
	a = a.To4()
	if a == nil {
		return nil, errors.New("invalid IPv4 address")
	}
	b := make(Attribute, len(a))
	copy(b, Attribute(a))
	return b, nil
}

// UserPassword decrypts the given  "User-Password"-encrypted (as defined in RFC
// 2865) Attribute, and returns the plaintext. An error is returned if the
// attribute length is invalid, the secret is empty, or the requestAuthenticator
// length is invalid.
func UserPassword(a Attribute, secret, requestAuthenticator []byte) ([]byte, error) {
	if len(a) < 16 || len(a) > 128 {
		return nil, errors.New("invalid attribute length (" + strconv.Itoa(len(a)) + ")")
	}
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(requestAuthenticator) != 16 {
		return nil, errors.New("invalid requestAuthenticator length (" + strconv.Itoa(len(requestAuthenticator)) + ")")
	}

	dec := make([]byte, 0, len(a))

	hash := md5.New()
	hash.Write(secret)
	hash.Write(requestAuthenticator)
	dec = hash.Sum(dec)

	for i, b := range a[:16] {
		dec[i] ^= b
	}

	for i := 16; i < len(a); i += 16 {
		hash.Reset()
		hash.Write(secret)
		hash.Write(a[i-16 : i])
		dec = hash.Sum(dec)

		for j, b := range a[i : i+16] {
			dec[i+j] ^= b
		}
	}

	if i := bytes.IndexByte(dec, 0); i > -1 {
		return dec[:i], nil
	}
	return dec, nil
}

// NewUserPassword returns a new "User-Password"-encrypted attribute from the
// given plaintext, secret, and requestAuthenticator. An error is returned if
// the plaintext is too long, the secret is empty, or the requestAuthenticator
// is an invalid length.
func NewUserPassword(plaintext, secret, requestAuthenticator []byte) (Attribute, error) {
	if len(plaintext) > 128 {
		return nil, errors.New("plaintext longer than 128 characters")
	}
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(requestAuthenticator) != 16 {
		return nil, errors.New("requestAuthenticator not 16-bytes")
	}

	chunks := len(plaintext) >> 4
	if chunks == 0 {
		chunks = 1
	}

	enc := make([]byte, 0, chunks*16)

	hash := md5.New()
	hash.Write(secret)
	hash.Write(requestAuthenticator)
	enc = hash.Sum(enc)

	for i, b := range plaintext[:16] {
		enc[i] ^= b
	}

	for i := 16; i < len(plaintext); i += 16 {
		hash.Reset()
		hash.Write(secret)
		hash.Write(enc[i-16 : i])
		enc = hash.Sum(enc)

		for j, b := range plaintext[i : i+16] {
			enc[i+j] ^= b
		}
	}

	return enc, nil
}

// Date returns the given Attribute as time.Time. An error is returned if the
// attribute is not 4 bytes long.
func Date(a Attribute) (time.Time, error) {
	if len(a) != 4 {
		return time.Time{}, errors.New("invalid length")
	}
	sec := binary.BigEndian.Uint32([]byte(a))
	return time.Unix(int64(sec), 0), nil
}

// NewDate returns a new Attribute from the given time.Time.
func NewDate(t time.Time) Attribute {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(t.Unix()))
	return a
}

// VendorSpecific returns the vendor ID and value from the given attribute. An
// error is returned if the attribute is less than 5 bytes long.
func VendorSpecific(a Attribute) (vendorID uint32, value Attribute, err error) {
	if len(a) < 5 {
		err = errors.New("invalid length")
		return
	}
	vendorID = binary.BigEndian.Uint32(a[:4])
	value = make([]byte, len(a)-4)
	copy(value, a[4:])
	return
}

// NewVendorSpecific returns a new vendor specific attribute with the given
// vendor ID and value.
func NewVendorSpecific(vendorID uint32, value Attribute) Attribute {
	a := make([]byte, 4+len(value))
	binary.BigEndian.PutUint32(a, vendorID)
	copy(a[4:], value)
	return a
}

// TODO: ipv6addr
// TODO: ipv6prefix
// TODO: ifid
// TODO: integer64
