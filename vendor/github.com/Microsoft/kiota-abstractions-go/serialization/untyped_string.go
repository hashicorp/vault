package serialization

// UntypedString defines an untyped string object.
type UntypedString struct {
	UntypedNode
}

// GetValue returns the string object.
func (un *UntypedString) GetValue() *string {
	castValue, ok := un.value.(*string)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedString creates a new UntypedString object.
func NewUntypedString(stringValue string) *UntypedString {
	m := &UntypedString{}
	m.value = &stringValue
	return m
}
