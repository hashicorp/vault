package radius

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"math"
	"net"
	"strconv"
	"time"
)

// ErrNoAttribute is returned when an attribute was not found when one was
// expected.
var ErrNoAttribute = errors.New("radius: attribute not found")

// Attribute is a wire encoded RADIUS attribute value.
type Attribute []byte

// Integer returns the given attribute as an integer. An error is returned if
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
	return v
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
	copy(b, a)
	return b
}

// NewBytes returns a new Attribute from the given byte slice. An error is
// returned if the slice is longer than 253.
func NewBytes(b []byte) (Attribute, error) {
	if len(b) > 253 {
		return nil, errors.New("value too long")
	}
	a := make(Attribute, len(b))
	copy(a, b)
	return a, nil
}

// IPAddr returns the given Attribute as an IPv4 IP address. An error is
// returned if the attribute is not 4 bytes long.
func IPAddr(a Attribute) (net.IP, error) {
	if len(a) != net.IPv4len {
		return nil, errors.New("invalid length")
	}
	b := make([]byte, net.IPv4len)
	copy(b, a)
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
	copy(b, a)
	return b, nil
}

// IPv6Addr returns the given Attribute as an IPv6 IP address. An error is
// returned if the attribute is not 16 bytes long.
func IPv6Addr(a Attribute) (net.IP, error) {
	if len(a) != net.IPv6len {
		return nil, errors.New("invalid length")
	}
	b := make([]byte, net.IPv6len)
	copy(b, a)
	return b, nil
}

// NewIPv6Addr returns a new Attribute from the given IP address. An error is
// returned if the given address is not an IPv6 address.
func NewIPv6Addr(a net.IP) (Attribute, error) {
	a = a.To16()
	if a == nil {
		return nil, errors.New("invalid IPv6 address")
	}
	b := make(Attribute, len(a))
	copy(b, a)
	return b, nil
}

// IFID returns the given attribute as a 8-byte hardware address. An error is
// return if the attribute is not 8 bytes long.
func IFID(a Attribute) (net.HardwareAddr, error) {
	if len(a) != 8 {
		return nil, errors.New("invalid length")
	}
	ifid := make(net.HardwareAddr, len(a))
	copy(ifid, a)
	return ifid, nil
}

// NewIFID returns a new Attribute from the given hardware address. An error
// is returned if the address is not 8 bytes long.
func NewIFID(addr net.HardwareAddr) (Attribute, error) {
	if len(addr) != 8 {
		return nil, errors.New("invalid length")
	}
	attr := make(Attribute, len(addr))
	copy(attr, addr)
	return attr, nil
}

// UserPassword decrypts the given  "User-Password"-encrypted (as defined in RFC
// 2865) Attribute, and returns the plaintext. An error is returned if the
// attribute length is invalid, the secret is empty, or the requestAuthenticator
// length is invalid.
func UserPassword(a Attribute, secret, requestAuthenticator []byte) ([]byte, error) {
	if len(a) < 16 || len(a) > 128 || len(a)%16 != 0 {
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

	chunks := (len(plaintext) + 16 - 1) / 16
	if chunks == 0 {
		chunks = 1
	}

	enc := make([]byte, 0, chunks*16)

	hash := md5.New()
	hash.Write(secret)
	hash.Write(requestAuthenticator)
	enc = hash.Sum(enc)

	for i := 0; i < 16 && i < len(plaintext); i++ {
		enc[i] ^= plaintext[i]
	}

	for i := 16; i < len(plaintext); i += 16 {
		hash.Reset()
		hash.Write(secret)
		hash.Write(enc[i-16 : i])
		enc = hash.Sum(enc)

		for j := 0; j < 16 && i+j < len(plaintext); j++ {
			enc[i+j] ^= plaintext[i+j]
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
	sec := binary.BigEndian.Uint32(a)
	return time.Unix(int64(sec), 0), nil
}

// NewDate returns a new Attribute from the given time.Time.
func NewDate(t time.Time) (Attribute, error) {
	unix := t.Unix()
	if unix > math.MaxUint32 {
		return nil, errors.New("time out of range")
	}
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(t.Unix()))
	return a, nil
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
func NewVendorSpecific(vendorID uint32, value Attribute) (Attribute, error) {
	if len(value) > 249 {
		return nil, errors.New("value too long")
	}
	a := make([]byte, 4+len(value))
	binary.BigEndian.PutUint32(a, vendorID)
	copy(a[4:], value)
	return a, nil
}

// Integer64 returns the given attribute as an integer. An error is returned if
// the attribute is not 8 bytes long.
func Integer64(a Attribute) (uint64, error) {
	if len(a) != 8 {
		return 0, errors.New("invalid length")
	}
	return binary.BigEndian.Uint64(a), nil
}

// NewInteger64 creates a new Attribute from the given integer value.
func NewInteger64(i uint64) Attribute {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(v, i)
	return v
}

// Short returns the given attribute as an integer. An error is returned if
// the attribute is not 2 bytes long.
func Short(a Attribute) (uint16, error) {
	if len(a) != 2 {
		return 0, errors.New("invalid length")
	}
	return binary.BigEndian.Uint16(a), nil
}

// NewShort creates a new Attribute from the given integer value.
func NewShort(i uint16) Attribute {
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(v, i)
	return v
}

// TLV returns a components of a Type-Length-Value (TLV) attribute.
func TLV(a Attribute) (tlvType byte, tlvValue Attribute, err error) {
	if len(a) < 3 || len(a) > 255 || int(a[1]) != len(a) {
		err = errors.New("invalid length")
		return
	}
	tlvType = a[0]
	tlvValue = make(Attribute, len(a)-2)
	copy(tlvValue, a[2:])
	return
}

// NewTLV returns a new TLV attribute.
func NewTLV(tlvType byte, tlvValue Attribute) (Attribute, error) {
	if len(tlvValue) < 1 || len(tlvValue) > 253 {
		return nil, errors.New("invalid value length")
	}
	a := make(Attribute, 1+1+len(tlvValue))
	a[0] = tlvType
	a[1] = byte(1 + 1 + len(tlvValue))
	copy(a[2:], tlvValue)
	return a, nil
}

// NewTunnelPassword returns an RFC 2868 encrypted Tunnel-Password.
// A tag must be added on to the returned Attribute.
func NewTunnelPassword(password, salt, secret, requestAuthenticator []byte) (Attribute, error) {
	if len(password) > 249 {
		return nil, errors.New("invalid password length")
	}
	if len(salt) != 2 {
		return nil, errors.New("invalid salt length")
	}
	if salt[0]&0x80 != 0x80 { // MSB must be 1
		return nil, errors.New("invalid salt")
	}
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(requestAuthenticator) != 16 {
		return nil, errors.New("invalid requestAuthenticator length")
	}

	chunks := (1 + len(password) + 16 - 1) / 16
	if chunks == 0 {
		chunks = 1
	}

	attr := make([]byte, 2+chunks*16)
	copy(attr[:2], salt)
	attr[2] = byte(len(password))
	copy(attr[3:], password)

	hash := md5.New()
	var b [md5.Size]byte

	for chunk := 0; chunk < chunks; chunk++ {
		hash.Reset()

		hash.Write(secret)
		if chunk == 0 {
			hash.Write(requestAuthenticator)
			hash.Write(salt)
		} else {
			hash.Write(attr[2+(chunk-1)*16 : 2+chunk*16])
		}
		hash.Sum(b[:0])

		for i := 0; i < 16; i++ {
			attr[2+chunk*16+i] ^= b[i]
		}
	}

	return attr, nil
}

// TunnelPassword decrypts an RFC 2868 encrypted Tunnel-Password.
// The Attribute must not be prefixed with a tag.
// The requestAuthenticator must be from the Access-Request packet.
func TunnelPassword(a Attribute, secret, requestAuthenticator []byte) (password, salt []byte, err error) {
	if len(a) > 252 || len(a) < 18 || (len(a)-2)%16 != 0 {
		err = errors.New("invalid length")
		return
	}
	if len(secret) == 0 {
		err = errors.New("empty secret")
		return
	}
	if len(requestAuthenticator) != 16 {
		err = errors.New("invalid requestAuthenticator length")
		return
	}
	if a[0]&0x80 != 0x80 { // salt MSB must be 1
		err = errors.New("invalid salt")
		return
	}

	salt = append([]byte(nil), a[:2]...)
	a = a[2:]

	chunks := len(a) / 16
	plaintext := make([]byte, chunks*16)

	hash := md5.New()
	var b [md5.Size]byte

	for chunk := 0; chunk < chunks; chunk++ {
		hash.Reset()

		hash.Write(secret)
		if chunk == 0 {
			hash.Write(requestAuthenticator)
			hash.Write(salt)
		} else {
			hash.Write(a[(chunk-1)*16 : chunk*16])
		}
		hash.Sum(b[:0])

		for i := 0; i < 16; i++ {
			plaintext[chunk*16+i] = a[chunk*16+i] ^ b[i]
		}
	}

	passwordLength := plaintext[0]
	if int(passwordLength) > (len(plaintext) - 1) {
		err = errors.New("invalid password length")
		return
	}
	password = plaintext[1 : 1+passwordLength]
	return
}

func NewIPv6Prefix(prefix *net.IPNet) (Attribute, error) {
	if prefix == nil {
		return nil, errors.New("nil prefix")
	}

	if len(prefix.IP) != net.IPv6len {
		return nil, errors.New("IP is not IPv6")
	}

	ones, bits := prefix.Mask.Size()
	if bits != net.IPv6len*8 {
		return nil, errors.New("mask is not IPv6")
	}

	attr := make(Attribute, 2+((ones+7)/8))
	// attr[0] = 0x00
	attr[1] = byte(ones)
	copy(attr[2:], prefix.IP)

	// clear final non-mask bits
	if i := uint(ones % 8); i != 0 {
		for ; i < 8; i++ {
			attr[len(attr)-1] &^= 1 << (7 - i)
		}
	}

	return attr, nil
}

func IPv6Prefix(a Attribute) (*net.IPNet, error) {
	if len(a) < 2 || len(a) > 18 {
		return nil, errors.New("invalid length")
	}

	prefixLength := int(a[1])
	if prefixLength > net.IPv6len*8 {
		return nil, errors.New("invalid prefix length")
	}

	ip := make(net.IP, net.IPv6len)
	copy(ip, a[2:])

	bit := uint(prefixLength % 8)
	for octet := prefixLength / 8; octet < len(ip); octet++ {
		for ; bit < 8; bit++ {
			if ip[octet]&(1<<(7-bit)) != 0 {
				return nil, errors.New("invalid prefix data")
			}
		}
		bit = 0
	}

	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(prefixLength, net.IPv6len*8),
	}, nil
}
