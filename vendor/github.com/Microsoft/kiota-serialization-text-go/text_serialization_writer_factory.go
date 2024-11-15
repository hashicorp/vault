package textserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TextSerializationWriterFactory implements SerializationWriterFactory for text.
type TextSerializationWriterFactory struct {
}

// NewTextSerializationWriterFactory creates a new instance of the TextSerializationWriterFactory.
func NewTextSerializationWriterFactory() *TextSerializationWriterFactory {
	return &TextSerializationWriterFactory{}
}

// GetValidContentType returns the valid content type for the SerializationWriterFactoryRegistry
func (f *TextSerializationWriterFactory) GetValidContentType() (string, error) {
	return "text/plain", nil
}

// GetSerializationWriter returns the relevant SerializationWriter instance for the given content type
func (f *TextSerializationWriterFactory) GetSerializationWriter(contentType string) (absser.SerializationWriter, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewTextSerializationWriter(), nil
	}
}
