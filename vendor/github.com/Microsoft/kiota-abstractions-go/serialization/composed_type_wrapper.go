package serialization

// ComposedTypeWrapper defines a contract for models that are wrappers for composed types.
type ComposedTypeWrapper interface {
	// GetIsComposedType returns true if the type is composed, false otherwise.
	GetIsComposedType() bool
}
