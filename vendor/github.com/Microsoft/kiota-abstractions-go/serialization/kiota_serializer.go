package serialization

import (
	"errors"
)

// Serialize serializes the given model into a byte array.
func Serialize(contentType string, model Parsable) ([]byte, error) {
	writer, err := getSerializationWriter(contentType, model)
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	err = writer.WriteObjectValue("", model)
	if err != nil {
		return nil, err
	}
	return writer.GetSerializedContent()
}

// SerializeCollection serializes the given models into a byte array.
func SerializeCollection(contentType string, models []Parsable) ([]byte, error) {
	writer, err := getSerializationWriter(contentType, models)
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	err = writer.WriteCollectionOfObjectValues("", models)
	if err != nil {
		return nil, err
	}
	return writer.GetSerializedContent()
}
func getSerializationWriter(contentType string, value interface{}) (SerializationWriter, error) {
	if contentType == "" {
		return nil, errors.New("the content type is empty")
	}
	if value == nil {
		return nil, errors.New("the value is empty")
	}
	writer, err := DefaultSerializationWriterFactoryInstance.GetSerializationWriter(contentType)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

// Deserialize deserializes the given byte array into a model.
func Deserialize(contentType string, content []byte, parsableFactory ParsableFactory) (Parsable, error) {
	node, err := getParseNode(contentType, content, parsableFactory)
	if err != nil {
		return nil, err
	}
	result, err := node.GetObjectValue(parsableFactory)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func getParseNode(contentType string, content []byte, parsableFactory ParsableFactory) (ParseNode, error) {
	if contentType == "" {
		return nil, errors.New("the content type is empty")
	}
	if content == nil || len(content) == 0 {
		return nil, errors.New("the content is empty")
	}
	if parsableFactory == nil {
		return nil, errors.New("the parsable factory is empty")
	}
	node, err := DefaultParseNodeFactoryInstance.GetRootParseNode(contentType, content)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// DeserializeCollection deserializes the given byte array into a collection of models.
func DeserializeCollection(contentType string, content []byte, parsableFactory ParsableFactory) ([]Parsable, error) {
	node, err := getParseNode(contentType, content, parsableFactory)
	if err != nil {
		return nil, err
	}
	result, err := node.GetCollectionOfObjectValues(parsableFactory)
	if err != nil {
		return nil, err
	}
	return result, nil
}
