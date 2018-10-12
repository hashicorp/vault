package pgtype

import (
	"database/sql/driver"
	"net"

	"github.com/pkg/errors"
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

	switch value := src.(type) {
	case net.IPNet:
		*dst = Inet{IPNet: &value, Status: Present}
	case *net.IPNet:
		*dst = Inet{IPNet: value, Status: Present}
	case net.IP:
		bitCount := len(value) * 8
		mask := net.CIDRMask(bitCount, bitCount)
		*dst = Inet{IPNet: &net.IPNet{Mask: mask, IP: value}, Status: Present}
	case string:
		_, ipnet, err := net.ParseCIDR(value)
		if err != nil {
			return err
		}
		*dst = Inet{IPNet: ipnet, Status: Present}
	default:
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Inet", value)
	}

	return nil
}

func (dst *Inet) Get() interface{} {
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
				return errors.Errorf("cannot assign %v to %T", src, dst)
			}
			*v = make(net.IP, len(src.IPNet.IP))
			copy(*v, src.IPNet.IP)
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return errors.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *Inet) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Inet{Status: Null}
		return nil
	}

	var ipnet *net.IPNet
	var err error

	if ip := net.ParseIP(string(src)); ip != nil {
		ipv4 := ip.To4()
		if ipv4 != nil {
			ip = ipv4
		}
		bitCount := len(ip) * 8
		mask := net.CIDRMask(bitCount, bitCount)
		ipnet = &net.IPNet{Mask: mask, IP: ip}
	} else {
		_, ipnet, err = net.ParseCIDR(string(src))
		if err != nil {
			return err
		}
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
		return errors.Errorf("Received an invalid size for a inet: %d", len(src))
	}

	// ignore family
	bits := src[1]
	// ignore is_cidr
	addressLength := src[3]

	var ipnet net.IPNet
	ipnet.IP = make(net.IP, int(addressLength))
	copy(ipnet.IP, src[4:])
	ipnet.Mask = net.CIDRMask(int(bits), int(addressLength)*8)

	*dst = Inet{IPNet: &ipnet, Status: Present}

	return nil
}

func (src *Inet) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.IPNet.String()...), nil
}

// EncodeBinary encodes src into w.
func (src *Inet) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
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
		return nil, errors.Errorf("Unexpected IP length: %v", len(src.IPNet.IP))
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

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Inet) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
