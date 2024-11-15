package radius

import (
	"errors"
)

// Type is the RADIUS attribute type.
type Type int

// TypeInvalid is a Type that can be used to represent an invalid RADIUS
// attribute type.
const TypeInvalid Type = -1

// AVP is an attribute-value pair.
// It contains an attribute type and its wire data.
type AVP struct {
	Type
	Attribute
}

// Attributes is a list of RADIUS attributes.
type Attributes []*AVP

// ParseAttributes parses the wire-encoded RADIUS attributes and returns a new
// Attributes value. An error is returned if the buffer is malformed.
func ParseAttributes(b []byte) (Attributes, error) {
	var attrs Attributes

	for len(b) > 0 {
		if len(b) < 2 {
			return nil, errors.New("short buffer")
		}
		length := int(b[1])
		if length > len(b) || length < 2 || length > 255 {
			return nil, errors.New("invalid attribute length")
		}

		avp := &AVP{
			Type: Type(b[0]),
		}
		if length > 2 {
			avp.Attribute = append(Attribute(nil), b[2:length]...)
		}
		attrs = append(attrs, avp)

		b = b[length:]
	}

	return attrs, nil
}

// Add appends the given Attribute to the list of attributes.
func (a *Attributes) Add(key Type, value Attribute) {
	*a = append(*a, &AVP{
		Type:      key,
		Attribute: value,
	})
}

// Del removes all Attributes of the given type from a.
func (a *Attributes) Del(key Type) {
	for i := 0; i < len(*a); {
		if (*a)[i].Type == key {
			*a = append((*a)[:i], (*a)[i+1:]...)
		} else {
			i++
		}
	}
}

// Get returns the first Attribute of Type key. nil is returned if no Attribute
// of Type key exists in a.
func (a *Attributes) Get(key Type) Attribute {
	attr, _ := a.Lookup(key)
	return attr
}

// Lookup returns the first Attribute of Type key. nil and false is returned if
// no Attribute of Type key exists in a.
func (a *Attributes) Lookup(key Type) (Attribute, bool) {
	for _, attr := range *a {
		if attr.Type == key {
			return attr.Attribute, true
		}
	}
	return nil, false
}

// Set removes all Attributes of Type key and appends value.
func (a *Attributes) Set(key Type, value Attribute) {
	foundKey := false
	for i := 0; i < len(*a); {
		if (*a)[i].Type == key {
			if foundKey {
				*a = append((*a)[:i], (*a)[i+1:]...)
			} else {
				(*a)[i] = &AVP{
					Type:      key,
					Attribute: value,
				}
				foundKey = true
				i++
			}
		} else {
			i++
		}
	}
	if !foundKey {
		a.Add(key, value)
	}
}

func (a Attributes) encodeTo(b []byte) {
	for _, attr := range a {
		if attr.Type < 0 || 255 < attr.Type || len(attr.Attribute) > 253 {
			continue
		}
		size := 1 + 1 + len(attr.Attribute)
		b[0] = byte(attr.Type)
		b[1] = byte(size)
		copy(b[2:], attr.Attribute)
		b = b[size:]
	}
}

// AttributesEncodedLen returns the encoded length of all attributes in a. An error is
// returned if any attribute in a exceeds the permitted size.
func AttributesEncodedLen(a Attributes) (int, error) {
	var n int
	for _, attr := range a {
		if attr.Type < 0 || 255 < attr.Type {
			continue
		}
		if len(attr.Attribute) > 253 {
			return 0, errors.New("radius: attribute too large")
		}
		n += 1 + 1 + len(attr.Attribute)
	}
	return n, nil
}
