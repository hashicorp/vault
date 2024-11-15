package serialization

type UntypedNodeable interface {
	Parsable
	GetIsUntypedNode() bool
}

// Base model for an untyped object.
type UntypedNode struct {
	value any
}

// GetIsUntypedNode returns true if the node is untyped, false otherwise.
func (m *UntypedNode) GetIsUntypedNode() bool {
	return true
}

// GetValue gets the underlying object value.
func (m *UntypedNode) GetValue() any {
	return m.value
}

// Serialize writes the objects properties to the current writer.
func (m *UntypedNode) Serialize(writer SerializationWriter) error {
	// Serialize the entity
	return nil
}

// GetFieldDeserializers returns the deserialization information for this object.
func (m *UntypedNode) GetFieldDeserializers() map[string]func(ParseNode) error {
	return make(map[string]func(ParseNode) error)
}

// NewUntypedNode instantiates a new untyped node and sets the default values.
func NewUntypedNode(value any) *UntypedNode {
	m := &UntypedNode{}
	m.value = value
	return m
}

// CreateUntypedNodeFromDiscriminatorValue a new untyped node and from a parse node.
func CreateUntypedNodeFromDiscriminatorValue(parseNode ParseNode) (Parsable, error) {
	return NewUntypedNode(nil), nil
}
