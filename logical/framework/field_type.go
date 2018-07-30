package framework

import (
	"net/http"
	"net/textproto"
)

// FieldType is the enum of types that a field can be.
type FieldType uint

const (
	TypeInvalid FieldType = 0
	TypeString  FieldType = iota
	TypeInt
	TypeBool
	TypeMap

	// TypeDurationSecond represent as seconds, this can be either an
	// integer or go duration format string (e.g. 24h)
	TypeDurationSecond

	// TypeSlice represents a slice of any type
	TypeSlice

	// TypeStringSlice is a helper for TypeSlice that returns a sanitized
	// slice of strings
	TypeStringSlice

	// TypeCommaStringSlice is a helper for TypeSlice that returns a sanitized
	// slice of strings and also supports parsing a comma-separated list in
	// a string field
	TypeCommaStringSlice

	// TypeLowerCaseString is a helper for TypeString that returns a lowercase
	// version of the provided string
	TypeLowerCaseString

	// TypeNameString represents a name that is URI safe and follows specific
	// rules.  These rules include start and end with an alphanumeric
	// character and characters in the middle can be alphanumeric or . or -.
	TypeNameString

	// TypeKVPairs allows you to represent the data as a map or a list of
	// equal sign delimited key pairs
	TypeKVPairs

	// TypeCommaIntSlice is a helper for TypeSlice that returns a sanitized
	// slice of Ints
	TypeCommaIntSlice

	// TypeHeader is a helper for sending request headers through to Vault.
	// For instance, the AWS and AliCloud credential plugins both act as a
	// benevolent MITM for a request, and the headers are sent through and
	// parsed.
	// IMPORTANT NOTES:
	//   - Under the hood, http.Header is a map[string][]string and its Get
	//     method will only return the first value. To retrieve all values,
	//     you must access them like you would in a map.
	//   - This implementation of http.Header has CASE-INSENSITIVE keys. All
	//     keys are converted to Title Case on their way in so that the Get
	//     method will match the keys no matter what case you use when
	//     looking for them.
	TypeHeader
)

func (t FieldType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeLowerCaseString:
		return "lowercase string"
	case TypeNameString:
		return "name string"
	case TypeInt:
		return "int"
	case TypeBool:
		return "bool"
	case TypeMap:
		return "map"
	case TypeKVPairs:
		return "keypair"
	case TypeDurationSecond:
		return "duration (sec)"
	case TypeSlice, TypeStringSlice, TypeCommaStringSlice, TypeCommaIntSlice:
		return "slice"
	case TypeHeader:
		return "header"
	default:
		return "unknown type"
	}
}

/*
	The reason we wrap an http.Header here is - under the hood,
	an http.Header is a map[string][]string and it does some magic
	with casing that can make it difficult to correctly iterate
	all header values.

	For example:
	h := http.Header{}
	h.Add("hello", "world")
	h.Add("hello", "monde")

	You'd expect that the header now have "hello" for a key, and
	both "world" and "monde" for a value. But it doesn't. "Hello"
	is now capitalized in the map, but "world" isn't.

	Later, when you do this:
	h.Get("hello")

	You'll receive only "world".

	If you try to solve for this by doing this:
	h["hello"]

	You'll receive nothing back, because remember, it's now in the
	map as "Hello.

	To avoid bugs like this, we provide one more method. GetAll,
	which returns all the values for a key.
*/
func NewHeader() Header {
	return Header{http.Header{}}
}

type Header struct {
	http.Header
}

func (h Header) GetAll(key string) []string {
	return h.Header[textproto.CanonicalMIMEHeaderKey(key)]
}
