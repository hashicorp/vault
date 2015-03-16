package framework

//go:generate stringer -type=FieldType field_type.go

// FieldType is the enum of types that a field can be.
type FieldType uint

const (
	TypeInvalid FieldType = 0
	TypeString  FieldType = iota
	TypeInt
	TypeBool
)

// FieldType has more methods defined on it in backend.go. They aren't
// in this file since stringer doesn't like that.
