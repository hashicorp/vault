package serialization

// SerializationWriterFactory defines the contract for a factory that creates SerializationWriter instances.
type SerializationWriterFactory interface {
	// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
	GetValidContentType() (string, error)
	// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
	GetSerializationWriter(contentType string) (SerializationWriter, error)
}
