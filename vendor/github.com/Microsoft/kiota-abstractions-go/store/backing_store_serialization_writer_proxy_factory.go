package store

import (
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

// BackingStoreSerializationWriterProxyFactory Backing Store implementation for serialization.SerializationWriterFactory
type BackingStoreSerializationWriterProxyFactory struct {
	factory serialization.SerializationWriterFactory
}

func (b *BackingStoreSerializationWriterProxyFactory) GetValidContentType() (string, error) {
	return b.factory.GetValidContentType()
}

func (b *BackingStoreSerializationWriterProxyFactory) GetSerializationWriter(contentType string) (serialization.SerializationWriter, error) {
	return b.factory.GetSerializationWriter(contentType)
}

// NewBackingStoreSerializationWriterProxyFactory Initializes a new instance of BackingStoreSerializationWriterProxyFactory
func NewBackingStoreSerializationWriterProxyFactory(factory serialization.SerializationWriterFactory) *BackingStoreSerializationWriterProxyFactory {
	proxyFactory := serialization.NewSerializationWriterProxyFactory(factory, func(parsable serialization.Parsable) error {
		if backedModel, ok := parsable.(BackedModel); ok && backedModel.GetBackingStore() != nil {
			backedModel.GetBackingStore().SetReturnOnlyChangedValues(true)
		}
		return nil
	}, func(parsable serialization.Parsable) error {
		if backedModel, ok := parsable.(BackedModel); ok && backedModel.GetBackingStore() != nil {
			store := backedModel.GetBackingStore()
			store.SetReturnOnlyChangedValues(false)
			store.SetInitializationCompleted(true)
		}
		return nil
	}, func(parsable serialization.Parsable, writer serialization.SerializationWriter) error {
		if backedModel, ok := parsable.(BackedModel); ok && backedModel.GetBackingStore() != nil {

			nilValues := backedModel.GetBackingStore().EnumerateKeysForValuesChangedToNil()
			for _, k := range nilValues {
				err := writer.WriteNullValue(k)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return &BackingStoreSerializationWriterProxyFactory{
		factory: proxyFactory,
	}
}
