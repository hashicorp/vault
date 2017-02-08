package radius // import "layeh.com/radius"

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
	"unicode/utf8"
)

// The base attribute value formats that are defined in RFC 2865.
var (
	// string
	AttributeText AttributeCodec
	// []byte
	AttributeString AttributeCodec
	// net.IP
	AttributeAddress AttributeCodec
	// uint32
	AttributeInteger AttributeCodec
	// time.Time
	AttributeTime AttributeCodec
	// []byte
	AttributeUnknown AttributeCodec
)

type attributeText struct{}

func (attributeText) Decode(packet *Packet, value []byte) (interface{}, error) {
	if !utf8.Valid(value) {
		return nil, errors.New("radius: text attribute is not valid UTF-8")
	}
	return string(value), nil
}

func (attributeText) Encode(packet *Packet, value interface{}) ([]byte, error) {
	str, ok := value.(string)
	if ok {
		return []byte(str), nil
	}
	raw, ok := value.([]byte)
	if ok {
		return raw, nil
	}
	return nil, errors.New("radius: text attribute must be string or []byte")
}

type attributeString struct{}

func (attributeString) Decode(packet *Packet, value []byte) (interface{}, error) {
	v := make([]byte, len(value))
	copy(v, value)
	return v, nil
}

func (attributeString) Encode(packet *Packet, value interface{}) ([]byte, error) {
	raw, ok := value.([]byte)
	if ok {
		return raw, nil
	}
	str, ok := value.(string)
	if ok {
		return []byte(str), nil
	}
	return nil, errors.New("radius: string attribute must be []byte or string")
}

type attributeAddress struct{}

func (attributeAddress) Decode(packet *Packet, value []byte) (interface{}, error) {
	if len(value) != net.IPv4len {
		return nil, errors.New("radius: address attribute has invalid size")
	}
	v := make([]byte, len(value))
	copy(v, value)
	return net.IP(v), nil
}

func (attributeAddress) Encode(packet *Packet, value interface{}) ([]byte, error) {
	ip, ok := value.(net.IP)
	if !ok {
		return nil, errors.New("radius: address attribute must be net.IP")
	}
	ip = ip.To4()
	if ip == nil {
		return nil, errors.New("radius: address attribute must be an IPv4 net.IP")
	}
	return []byte(ip), nil
}

type attributeInteger struct{}

func (attributeInteger) Decode(packet *Packet, value []byte) (interface{}, error) {
	if len(value) != 4 {
		return nil, errors.New("radius: integer attribute has invalid size")
	}
	return binary.BigEndian.Uint32(value), nil
}

func (attributeInteger) Encode(packet *Packet, value interface{}) ([]byte, error) {
	integer, ok := value.(uint32)
	if !ok {
		return nil, errors.New("radius: integer attribute must be uint32")
	}
	raw := make([]byte, 4)
	binary.BigEndian.PutUint32(raw, integer)
	return raw, nil
}

type attributeTime struct{}

func (attributeTime) Decode(packet *Packet, value []byte) (interface{}, error) {
	if len(value) != 4 {
		return nil, errors.New("radius: time attribute has invalid size")
	}
	return time.Unix(int64(binary.BigEndian.Uint32(value)), 0), nil
}

func (attributeTime) Encode(packet *Packet, value interface{}) ([]byte, error) {
	timestamp, ok := value.(time.Time)
	if !ok {
		return nil, errors.New("radius: time attribute must be time.Time")
	}
	raw := make([]byte, 4)
	binary.BigEndian.PutUint32(raw, uint32(timestamp.Unix()))
	return raw, nil
}
