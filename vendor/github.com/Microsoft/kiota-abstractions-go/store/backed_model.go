package store

// BackedModel defines methods for models that support a backing store
type BackedModel interface {
	// GetBackingStore returns the store that is backing the model.
	GetBackingStore() BackingStore
}
