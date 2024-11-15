package jsonserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// JsonSerializationWriterFactory implements SerializationWriterFactory for JSON.
type JsonSerializationWriterFactory struct {
}

// NewJsonSerializationWriterFactory creates a new instance of the JsonSerializationWriterFactory.
func NewJsonSerializationWriterFactory() *JsonSerializationWriterFactory {
	return &JsonSerializationWriterFactory{}
}

// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
func (f *JsonSerializationWriterFactory) GetValidContentType() (string, error) {
	return "application/json", nil
}

// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
func (f *JsonSerializationWriterFactory) GetSerializationWriter(contentType string) (absser.SerializationWriter, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewJsonSerializationWriter(), nil
	}
}
