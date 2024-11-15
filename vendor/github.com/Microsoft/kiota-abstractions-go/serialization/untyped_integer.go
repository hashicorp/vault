package serialization

// UntypedInteger defines an untyped integer value.
type UntypedInteger struct {
	UntypedNode
}

// GetValue returns the int32 value.
func (un *UntypedInteger) GetValue() *int32 {
	castValue, ok := un.value.(*int32)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedInteger creates a new UntypedInteger object.
func NewUntypedInteger(int32Value int32) *UntypedInteger {
	m := &UntypedInteger{}
	m.value = &int32Value
	return m
}
