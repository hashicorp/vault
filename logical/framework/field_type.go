package framework

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
)

func (t FieldType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeInt:
		return "int"
	case TypeBool:
		return "bool"
	case TypeMap:
		return "map"
	case TypeDurationSecond:
		return "duration (sec)"
	case TypeSlice, TypeStringSlice, TypeCommaStringSlice:
		return "slice"
	default:
		return "unknown type"
	}
}
