package serialization

// UntypedBoolean defines an untyped boolean object.
type UntypedBoolean struct {
	UntypedNode
}

// GetValue returns the bool value.
func (un *UntypedBoolean) GetValue() *bool {
	castValue, ok := un.value.(*bool)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedBoolean creates a new UntypedBoolean object.
func NewUntypedBoolean(boolValue bool) *UntypedBoolean {
	m := &UntypedBoolean{}
	m.value = &boolValue
	return m
}
