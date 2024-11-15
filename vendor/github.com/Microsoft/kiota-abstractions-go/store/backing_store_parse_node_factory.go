package store

import "github.com/microsoft/kiota-abstractions-go/serialization"

// BackingStoreParseNodeFactory Backing Store implementation for serialization.ParseNodeFactory
type BackingStoreParseNodeFactory struct {
	serialization.ParseNodeFactory
}

// NewBackingStoreParseNodeFactory Initializes a new instance of BackingStoreParseNodeFactory
func NewBackingStoreParseNodeFactory(factory serialization.ParseNodeFactory) *BackingStoreParseNodeFactory {
	proxyFactory := serialization.NewParseNodeProxyFactory(factory, func(parsable serialization.Parsable) error {
		if backedModel, ok := parsable.(BackedModel); ok && backedModel.GetBackingStore() != nil {
			backedModel.GetBackingStore().SetInitializationCompleted(false)
		}
		return nil
	}, func(parsable serialization.Parsable) error {
		if backedModel, ok := parsable.(BackedModel); ok && backedModel.GetBackingStore() != nil {
			backedModel.GetBackingStore().SetInitializationCompleted(true)
		}
		return nil
	})

	return &BackingStoreParseNodeFactory{proxyFactory}
}
