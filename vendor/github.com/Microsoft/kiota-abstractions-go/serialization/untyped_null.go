package serialization

// UntypedNull defines a untyped nil object.
type UntypedNull struct {
	UntypedNode
}

// GetValue returns a nil value.
func (un *UntypedNull) GetValue() any {
	return nil
}

// NewUntypedString creates a new UntypedNull object.
func NewUntypedNull() *UntypedNull {
	m := &UntypedNull{}
	return m
}
