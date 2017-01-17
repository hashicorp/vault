package radius // import "layeh.com/radius"

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math"
)

type rfc2865UserPassword struct{}

func (rfc2865UserPassword) Decode(p *Packet, value []byte) (interface{}, error) {
	if p.Secret == nil {
		return nil, errors.New("radius: User-Password attribute requires Packet.Secret")
	}
	if len(value) < 16 || len(value) > 128 {
		return nil, errors.New("radius: invalid User-Password attribute length")
	}

	dec := make([]byte, 0, len(value))

	hash := md5.New()
	hash.Write(p.Secret)
	hash.Write(p.Authenticator[:])
	dec = hash.Sum(dec)

	for i, b := range value[:16] {
		dec[i] ^= b
	}

	for i := 16; i < len(value); i += 16 {
		hash.Reset()
		hash.Write(p.Secret)
		hash.Write(value[i-16 : i])
		dec = hash.Sum(dec)

		for j, b := range value[i : i+16] {
			dec[i+j] ^= b
		}
	}

	if i := bytes.IndexByte(dec, 0); i > -1 {
		return string(dec[:i]), nil
	}
	return string(dec), nil
}

func (rfc2865UserPassword) Encode(p *Packet, value interface{}) ([]byte, error) {
	if p.Secret == nil {
		return nil, errors.New("radius: User-Password attribute requires Packet.Secret")
	}
	var password []byte
	if bytePassword, ok := value.([]byte); !ok {
		strPassword, ok := value.(string)
		if !ok {
			return nil, errors.New("radius: User-Password attribute must be string or []byte")
		}
		password = []byte(strPassword)
	} else {
		password = bytePassword
	}

	if len(password) > 128 {
		return nil, errors.New("radius: User-Password longer than 128 characters")
	}

	chunks := int(math.Ceil(float64(len(password)) / 16.))
	if chunks == 0 {
		chunks = 1
	}

	enc := make([]byte, 0, chunks*16)

	hash := md5.New()
	hash.Write(p.Secret)
	hash.Write(p.Authenticator[:])
	enc = hash.Sum(enc)

	for i, b := range password[:16] {
		enc[i] ^= b
	}

	for i := 16; i < len(password); i += 16 {
		hash.Reset()
		hash.Write(p.Secret)
		hash.Write(enc[i-16 : i])
		enc = hash.Sum(enc)

		for j, b := range password[i : i+16] {
			enc[i+j] ^= b
		}
	}

	return enc, nil
}

// VendorSpecific defines RFC 2865's Vendor-Specific attribute.
type VendorSpecific struct {
	VendorID uint32
	Data     []byte
}

type rfc2865VendorSpecific struct{}

func (rfc2865VendorSpecific) Decode(p *Packet, value []byte) (interface{}, error) {
	if len(value) < 5 {
		return nil, errors.New("radius: Vendor-Specific attribute too small")
	}
	var attr VendorSpecific
	attr.VendorID = binary.BigEndian.Uint32(value[:4])
	attr.Data = make([]byte, len(value)-4)
	copy(attr.Data, value[4:])
	return attr, nil
}

func (rfc2865VendorSpecific) Encode(p *Packet, value interface{}) ([]byte, error) {
	attr, ok := value.(VendorSpecific)
	if !ok {
		return nil, errors.New("radius: Vendor-Specific attribute is not type VendorSpecific")
	}
	b := make([]byte, 4+len(attr.Data))
	binary.BigEndian.PutUint32(b[:4], attr.VendorID)
	copy(b[4:], attr.Data)
	return b, nil
}
