package serialization

// UntypedFloat defines an untyped float32 value.
type UntypedFloat struct {
	UntypedNode
}

// GetValue returns the float32 value.
func (un *UntypedFloat) GetValue() *float32 {
	castValue, ok := un.value.(*float32)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedFloat creates a new UntypedFloat object.
func NewUntypedFloat(float32Value float32) *UntypedFloat {
	m := &UntypedFloat{}
	m.value = &float32Value
	return m
}
