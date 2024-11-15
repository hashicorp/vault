package serialization

// UntypedDouble defines an untyped float64 object.
type UntypedDouble struct {
	UntypedNode
}

// GetValue returns the float64 value.
func (un *UntypedDouble) GetValue() *float64 {
	castValue, ok := un.value.(*float64)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedDouble creates a new UntypedDouble object.
func NewUntypedDouble(float64Value float64) *UntypedDouble {
	m := &UntypedDouble{}
	m.value = &float64Value
	return m
}
