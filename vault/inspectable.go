package vault

type Inspectable interface {
	// Returns a record view of a particular subsystem
	GetRecords(tag string) ([]map[string]any, error)
}

type Deserializable interface {
	// Converts a structure into a consummable map
	Deserialize() map[string]any
}
