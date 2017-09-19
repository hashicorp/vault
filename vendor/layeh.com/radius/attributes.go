package radius

import (
	"errors"
)

// Type is the RADIUS attribute type.
type Type int

// TypeInvalid is a Type that can be used to represent an invalid RADIUS
// attribute type.
const TypeInvalid Type = -1

// Attributes is a map of RADIUS attribute types to slice of Attributes.
type Attributes map[Type][]Attribute

// ParseAttributes parses the wire-encoded RADIUS attributes and returns a new
// Attributes value. An error is returned if the buffer is malformed.
func ParseAttributes(b []byte) (Attributes, error) {
	attrs := make(map[Type][]Attribute)

	for len(b) > 0 {
		if len(b) < 2 {
			return nil, errors.New("short buffer")
		}
		length := int(b[1])
		if length > len(b) || length > 253 {
			return nil, errors.New("invalid attribute length")
		}

		typ := Type(b[0])
		var value Attribute
		if length > 2 {
			value = make(Attribute, length-2)
			copy(value, b[2:])
		}
		attrs[typ] = append(attrs[typ], value)

		b = b[length:]
	}

	return attrs, nil
}

// Add appends the given Attribute to the map entry of the given type.
func (a Attributes) Add(key Type, value Attribute) {
	a[key] = append(a[key], value)
}

// Del removes all Attributes of the given type from a.
func (a Attributes) Del(key Type) {
	delete(a, key)
}

// Get returns the first Attribute of Type key. nil is returned if no Attribute
// of Type key exists in a.
func (a Attributes) Get(key Type) Attribute {
	attr, _ := a.Lookup(key)
	return attr
}

// Lookup returns the first Attribute of Type key. nil and false is returned if
// no Attribute of Type key exists in a.
func (a Attributes) Lookup(key Type) (Attribute, bool) {
	m := a[key]
	if len(m) == 0 {
		return nil, false
	}
	return m[0], true
}

// Set removes all Attributes of Type key and appends value.
func (a Attributes) Set(key Type, value Attribute) {
	a[key] = []Attribute{value}
}

// Len returns the total number of Attributes in a.
func (a Attributes) Len() int {
	var i int
	for _, s := range a {
		i += len(s)
	}
	return i
}

func (a Attributes) encodeTo(b []byte) {
	for typ, attrs := range a {
		if typ < 0 || typ > 255 {
			continue
		}
		for _, attr := range attrs {
			size := 2 + len(attr)
			b[0] = byte(typ)
			b[1] = byte(size)
			copy(b[2:], attr)
			b = b[size:]
		}
	}
}

func (a Attributes) wireSize() (bytes int) {
	for typ, attrs := range a {
		if typ < 0 || typ > 255 {
			continue
		}
		for _, attr := range attrs {
			// type field + length field + value field
			bytes += 1 + 1 + len(attr)
		}
	}
	return
}
