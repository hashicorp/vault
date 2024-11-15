package pgtype

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"net"
	"strings"
)

// Network address family is dependent on server socket.h value for AF_INET.
// In practice, all platforms appear to have the same value. See
// src/include/utils/inet.h for more information.
const (
	defaultAFInet  = 2
	defaultAFInet6 = 3
)

// Inet represents both inet and cidr PostgreSQL types.
type Inet struct {
	IPNet  *net.IPNet
	Status Status
}

func (dst *Inet) Set(src interface{}) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case net.IPNet:
		*dst = Inet{IPNet: &value, Status: Present}
	case net.IP:
		if len(value) == 0 {
			*dst = Inet{Status: Null}
		} else {
			bitCount := len(value) * 8
			mask := net.CIDRMask(bitCount, bitCount)
			*dst = Inet{IPNet: &net.IPNet{Mask: mask, IP: value}, Status: Present}
		}
	case string:
		ip, ipnet, err := net.ParseCIDR(value)
		if err != nil {
			ip := net.ParseIP(value)
			if ip == nil {
				return fmt.Errorf("unable to parse inet address: %s", value)
			}

			if ipv4 := maybeGetIPv4(value, ip); ipv4 != nil {
				ipnet = &net.IPNet{IP: ipv4, Mask: net.CIDRMask(32, 32)}
			} else {
				ipnet = &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
			}
		} else {
			ipnet.IP = ip
			if ipv4 := maybeGetIPv4(value, ipnet.IP); ipv4 != nil {
				ipnet.IP = ipv4
				if len(ipnet.Mask) == 16 {
					ipnet.Mask = ipnet.Mask[12:] // Not sure this is ever needed.
				}
			}
		}

		*dst = Inet{IPNet: ipnet, Status: Present}
	case *net.IPNet:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *net.IP:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Inet{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if tv, ok := src.(encoding.TextMarshaler); ok {
			text, err := tv.MarshalText()
			if err != nil {
				return fmt.Errorf("cannot marshal %v: %w", value, err)
			}
			return dst.Set(string(text))
		}
		if sv, ok := src.(fmt.Stringer); ok {
			return dst.Set(sv.String())
		}
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Inet", value)
	}

	return nil
}

// Convert the net.IP to IPv4, if appropriate.
//
// When parsing a string to a net.IP using net.ParseIP() and the like, we get a
// 16 byte slice for IPv4 addresses as well as IPv6 addresses. This function
// calls To4() to convert them to a 4 byte slice. This is useful as it allows
// users of the net.IP check for IPv4 addresses based on the length and makes
// it clear we are handling IPv4 as opposed to IPv6 or IPv4-mapped IPv6
// addresses.
func maybeGetIPv4(input string, ip net.IP) net.IP {
	// Do not do this if the provided input looks like IPv6. This is because
	// To4() on IPv4-mapped IPv6 addresses converts them to IPv4, which behave
	// different in some cases.
	if strings.Contains(input, ":") {
		return nil
	}

	return ip.To4()
}

func (dst Inet) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.IPNet
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Inet) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *net.IPNet:
			*v = net.IPNet{
				IP:   make(net.IP, len(src.IPNet.IP)),
				Mask: make(net.IPMask, len(src.IPNet.Mask)),
			}
			copy(v.IP, src.IPNet.IP)
			copy(v.Mask, src.IPNet.Mask)
			return nil
		case *net.IP:
			if oneCount, bitCount := src.IPNet.Mask.Size(); oneCount != bitCount {
				return fmt.Errorf("cannot assign %v to %T", src, dst)
			}
			*v = make(net.IP, len(src.IPNet.IP))
			copy(*v, src.IPNet.IP)
			return nil
		default:
			if tv, ok := dst.(encoding.TextUnmarshaler); ok {
				if err := tv.UnmarshalText([]byte(src.IPNet.String())); err != nil {
					return fmt.Errorf("cannot unmarshal %v to %T: %w", src, dst, err)
				}
				return nil
			}
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return fmt.Errorf("unable to assign to %T", dst)
		}
	case Null:
		return NullAssignTo(dst)
	}

	return fmt.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *Inet) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	var ipnet *net.IPNet
	var err error

	if ip := net.ParseIP(string(src)); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
		bitCount := len(ip) * 8
		mask := net.CIDRMask(bitCount, bitCount)
		ipnet = &net.IPNet{Mask: mask, IP: ip}
	} else {
		ip, ipnet, err = net.ParseCIDR(string(src))
		if err != nil {
			return err
		}
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
		ones, _ := ipnet.Mask.Size()
		*ipnet = net.IPNet{IP: ip, Mask: net.CIDRMask(ones, len(ip)*8)}
	}

	*dst = Inet{IPNet: ipnet, Status: Present}
	return nil
}

func (dst *Inet) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	if len(src) != 8 && len(src) != 20 {
		return fmt.Errorf("Received an invalid size for an inet: %d", len(src))
	}

	// ignore family
	bits := src[1]
	// ignore is_cidr
	addressLength := src[3]

	var ipnet net.IPNet
	ipnet.IP = make(net.IP, int(addressLength))
	copy(ipnet.IP, src[4:])
	if ipv4 := ipnet.IP.To4(); ipv4 != nil {
		ipnet.IP = ipv4
	}
	ipnet.Mask = net.CIDRMask(int(bits), len(ipnet.IP)*8)

	*dst = Inet{IPNet: &ipnet, Status: Present}

	return nil
}

func (src Inet) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.IPNet.String()...), nil
}

// EncodeBinary encodes src into w.
func (src Inet) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var family byte
	switch len(src.IPNet.IP) {
	case net.IPv4len:
		family = defaultAFInet
	case net.IPv6len:
		family = defaultAFInet6
	default:
		return nil, fmt.Errorf("Unexpected IP length: %v", len(src.IPNet.IP))
	}

	buf = append(buf, family)

	ones, _ := src.IPNet.Mask.Size()
	buf = append(buf, byte(ones))

	// is_cidr is ignored on server
	buf = append(buf, 0)

	buf = append(buf, byte(len(src.IPNet.IP)))

	return append(buf, src.IPNet.IP...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Inet) Scan(src interface{}) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Inet) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
