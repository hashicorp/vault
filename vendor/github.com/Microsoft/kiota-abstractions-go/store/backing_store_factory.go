package store

// BackingStoreFactory represents a factory function for a backing store
// initializes a new backing store object
type BackingStoreFactory func() BackingStore

// BackingStoreFactoryInstance returns a backing store instance.
// if none exists an instance of InMemoryBackingStore is initialized and returned
var BackingStoreFactoryInstance BackingStoreFactory = NewInMemoryBackingStore
