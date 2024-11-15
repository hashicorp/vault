package multipartserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MultipartSerializationWriterFactory implements SerializationWriterFactory for URI Multipart encoded.
type MultipartSerializationWriterFactory struct {
}

// NewMultipartSerializationWriterFactory creates a new instance of the MultipartSerializationWriterFactory.
func NewMultipartSerializationWriterFactory() *MultipartSerializationWriterFactory {
	return &MultipartSerializationWriterFactory{}
}

// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
func (f *MultipartSerializationWriterFactory) GetValidContentType() (string, error) {
	return "multipart/form-data", nil
}

// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
func (f *MultipartSerializationWriterFactory) GetSerializationWriter(contentType string) (absser.SerializationWriter, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewMultipartSerializationWriter(), nil
	}
}
