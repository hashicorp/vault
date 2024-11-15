package textserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TextParseNodeFactory is a ParseNodeFactory implementation for text
type TextParseNodeFactory struct {
}

// Creates a new TextParseNodeFactory
func NewTextParseNodeFactory() *TextParseNodeFactory {
	return &TextParseNodeFactory{}
}

// GetValidContentType returns the content type this factory's parse nodes can deserialize.
func (f *TextParseNodeFactory) GetValidContentType() (string, error) {
	return "text/plain", nil
}

// GetRootParseNode return a new ParseNode instance that is the root of the content
func (f *TextParseNodeFactory) GetRootParseNode(contentType string, content []byte) (absser.ParseNode, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewTextParseNode(content)
	}
}
