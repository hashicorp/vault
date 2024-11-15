package store

// BackingStore Stores model information in a different location than the object properties.
// Implementations can provide dirty tracking capabilities, caching capabilities or integration with 3rd party stores.
type BackingStore interface {

	// Get return a value from the backing store based on its key.
	// Returns null if the value hasn't changed and "ReturnOnlyChangedValues" is true
	Get(key string) (interface{}, error)

	// Set or updates the stored value for the given key.
	// Will trigger subscriptions callbacks.
	Set(key string, value interface{}) error

	// Enumerate returns all the values stored in the backing store. Values will be filtered if "ReturnOnlyChangedValues" is true.
	Enumerate() map[string]interface{}

	// EnumerateKeysForValuesChangedToNil returns the keys for all values that changed to null
	EnumerateKeysForValuesChangedToNil() []string

	// Subscribe registers a listener to any data change happening.
	// returns a subscriptionId which cah be used to reference the current subscription
	Subscribe(callback BackingStoreSubscriber) string

	// SubscribeWithId registers a listener to any data change happening and assigns the given id
	SubscribeWithId(callback BackingStoreSubscriber, subscriptionId string) error

	// Unsubscribe Removes a subscription from the store based on its subscription id.
	Unsubscribe(subscriptionId string) error

	// Clear Removes the data stored in the backing store. Doesn't trigger any subscription.
	Clear()

	// GetInitializationCompleted Track's status of object during serialization and deserialization.
	// this property is used to initialize subscriber notifications
	GetInitializationCompleted() bool

	// SetInitializationCompleted sets whether the initialization of the object and/or
	// the initial deserialization has been completed to track whether objects have changed
	SetInitializationCompleted(val bool)

	// GetReturnOnlyChangedValues is a flag that defines whether subscriber notifications should be sent only when
	// data has been changed
	GetReturnOnlyChangedValues() bool

	// SetReturnOnlyChangedValues Sets whether to return only values that have changed
	// since the initialization of the object when calling the Get and Enumerate method
	SetReturnOnlyChangedValues(val bool)
}
