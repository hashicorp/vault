package jsonserialization

import (
	"encoding/json"
	"reflect"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Unmarshal parses JSON-encoded data using a ParsableFactory and stores it in the value pointed to by model. To
// enable dirty tracking and better performance, set the DefaultParseNodeFactoryInstance for "application/json".
func Unmarshal[T absser.Parsable](data []byte, model *T, parser absser.ParsableFactory) error {
	jpn, err := absser.DefaultParseNodeFactoryInstance.GetRootParseNode("application/json", data)
	if err != nil {
		if jpn, err = NewJsonParseNode(data); err != nil {
			return err
		}
	}

	v, err := jpn.GetObjectValue(parser)
	if err != nil {
		return err
	}

	if v != nil {
		*model = v.(T)
	} else {
		// hand off to the std library to set model to its zero value
		return json.Unmarshal(data, model)
	}

	return nil
}

// Marshal JSON-encodes a Parsable value. To enable dirty tracking and better performance, set the
// DefaultSerializationWriterFactoryInstance for "application/json".
func Marshal(v absser.Parsable) ([]byte, error) {
	if vRef := reflect.ValueOf(v); !vRef.IsValid() || vRef.IsNil() {
		return []byte("null"), nil
	}

	serializer, err := absser.DefaultSerializationWriterFactoryInstance.GetSerializationWriter("application/json")
	if err != nil {
		serializer = NewJsonSerializationWriter()
	}

	defer serializer.Close()

	if err := v.Serialize(serializer); err != nil {
		return nil, err
	}

	return serializer.GetSerializedContent()
}
