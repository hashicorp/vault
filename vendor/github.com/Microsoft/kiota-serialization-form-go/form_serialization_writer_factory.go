package formserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FormSerializationWriterFactory implements SerializationWriterFactory for URI form encoded.
type FormSerializationWriterFactory struct {
}

// NewFormSerializationWriterFactory creates a new instance of the FormSerializationWriterFactory.
func NewFormSerializationWriterFactory() *FormSerializationWriterFactory {
	return &FormSerializationWriterFactory{}
}

// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
func (f *FormSerializationWriterFactory) GetValidContentType() (string, error) {
	return "application/x-www-form-urlencoded", nil
}

// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
func (f *FormSerializationWriterFactory) GetSerializationWriter(contentType string) (absser.SerializationWriter, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewFormSerializationWriter(), nil
	}
}
