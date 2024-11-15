package serialization

var jsonContentType = "application/json"

// SerializeToJson serializes the given model to JSON
func SerializeToJson(model Parsable) ([]byte, error) {
	return Serialize(jsonContentType, model)
}

// SerializeCollectionToJson serializes the given models to JSON
func SerializeCollectionToJson(models []Parsable) ([]byte, error) {
	return SerializeCollection(jsonContentType, models)
}

// DeserializeFromJson deserializes the given JSON to a model
func DeserializeFromJson(content []byte, parsableFactory ParsableFactory) (Parsable, error) {
	return Deserialize(jsonContentType, content, parsableFactory)
}

// DeserializeCollectionFromJson deserializes the given JSON to a collection of models
func DeserializeCollectionFromJson(content []byte, parsableFactory ParsableFactory) ([]Parsable, error) {
	return DeserializeCollection(jsonContentType, content, parsableFactory)
}
