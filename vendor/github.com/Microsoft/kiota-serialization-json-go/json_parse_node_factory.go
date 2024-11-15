package jsonserialization

import (
	"errors"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// JsonParseNodeFactory is a ParseNodeFactory implementation for JSON
type JsonParseNodeFactory struct {
}

// NewJsonParseNodeFactory creates a new JsonParseNodeFactory
func NewJsonParseNodeFactory() *JsonParseNodeFactory {
	return &JsonParseNodeFactory{}
}

// GetValidContentType returns the content type this factory's parse nodes can deserialize.
func (f *JsonParseNodeFactory) GetValidContentType() (string, error) {
	return "application/json", nil
}

// GetRootParseNode return a new ParseNode instance that is the root of the content
func (f *JsonParseNodeFactory) GetRootParseNode(contentType string, content []byte) (absser.ParseNode, error) {
	validType, err := f.GetValidContentType()
	if err != nil {
		return nil, err
	} else if contentType == "" {
		return nil, errors.New("contentType is empty")
	} else if contentType != validType {
		return nil, errors.New("contentType is not valid")
	} else {
		return NewJsonParseNode(content)
	}
}
