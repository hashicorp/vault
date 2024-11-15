package serialization

// UntypedObject defines an untyped object.
type UntypedObject struct {
	UntypedNode
}

// GetValue gets a map of the properties of the object.
func (un *UntypedObject) GetValue() map[string]UntypedNodeable {
	castValue, ok := un.value.(map[string]UntypedNodeable)
	if ok {
		return castValue
	}
	return nil
}

// NewUntypedObject creates a new UntypedObject object.
func NewUntypedObject(properties map[string]UntypedNodeable) *UntypedObject {
	m := &UntypedObject{}
	m.value = properties
	return m
}
