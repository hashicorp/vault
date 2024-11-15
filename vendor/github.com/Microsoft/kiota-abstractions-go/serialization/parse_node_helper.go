package serialization

import "errors"

// MergeDeserializersForIntersectionWrapper merges the given fields deserializers for an intersection type into a single collection.
func MergeDeserializersForIntersectionWrapper(targets ...Parsable) (map[string]func(ParseNode) error, error) {
	if len(targets) == 0 {
		return nil, errors.New("no targets provided")
	}
	if len(targets) == 1 {
		return targets[0].GetFieldDeserializers(), nil
	}
	deserializers := make(map[string]func(ParseNode) error)
	for _, target := range targets {
		if target != nil {
			for key, deserializer := range target.GetFieldDeserializers() {
				if _, ok := deserializers[key]; !ok {
					deserializers[key] = deserializer
				}
			}
		}
	}
	return deserializers, nil
}
