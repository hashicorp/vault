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

	// TypeCommaStringSlice represents a slice either as native slice or
	// a comma-seperated string (value1,value2 => ["value1", "value2"])
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
	case TypeCommaStringSlice:
		return "slice"
	default:
		return "unknown type"
	}
}
