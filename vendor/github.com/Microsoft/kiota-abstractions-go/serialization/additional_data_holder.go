package serialization

// AdditionalDataHolder defines a contract for models that can hold additional data besides the described properties.
type AdditionalDataHolder interface {
	// GetAdditionalData returns additional data of the object that doesn't belong to a field.
	GetAdditionalData() map[string]interface{}
	// SetAdditionalData sets additional data of the object that doesn't belong to a field.
	SetAdditionalData(value map[string]interface{})
}
