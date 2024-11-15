package store

import (
	"errors"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// BackingStoreSubscriber is a function signature for a listener to any changes to a backing store
// It takes a `key` that represents the name of the object that's changed,
// `oldVal` is the previous value
// `newVal` is the newly assigned value
type BackingStoreSubscriber func(key string, oldVal interface{}, newVal interface{})

type InMemoryBackingStore struct {
	returnOnlyChangedValues bool
	initializationCompleted bool
	store                   map[string]interface{}
	subscribers             map[string]BackingStoreSubscriber
	changedValues           map[string]bool
}

// NewInMemoryBackingStore returns a new instance of an in memory backing store
// this function also provides an implementation of a BackingStoreFactory
func NewInMemoryBackingStore() BackingStore {
	return &InMemoryBackingStore{
		returnOnlyChangedValues: false,
		initializationCompleted: true,
		store:                   make(map[string]interface{}),
		subscribers:             make(map[string]BackingStoreSubscriber),
		changedValues:           make(map[string]bool),
	}
}

func (i *InMemoryBackingStore) Get(key string) (interface{}, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("key cannot be an empty string")
	}

	objectVal := i.store[key]

	if (i.GetReturnOnlyChangedValues() && i.changedValues[key]) || !i.GetReturnOnlyChangedValues() {
		return objectVal, nil
	} else {
		return nil, nil
	}
}

func (i *InMemoryBackingStore) Set(key string, value interface{}) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("key cannot be an empty string")
	}

	current := i.store[key]

	// check if objects values have changed
	if current == nil || hasChanged(current, value) {
		// track changed key
		i.changedValues[key] = i.GetInitializationCompleted()

		// update changed values
		i.store[key] = value

		// notify subs
		for _, subscriber := range i.subscribers {
			subscriber(key, current, value)
		}
	}

	return nil
}

func hasChanged(current interface{}, value interface{}) bool {
	kind := reflect.ValueOf(current).Kind()
	if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Struct {
		return !reflect.DeepEqual(current, value)
	} else {
		return current != value
	}
}

func (i *InMemoryBackingStore) Enumerate() map[string]interface{} {
	items := make(map[string]interface{})

	for k, v := range i.store {
		if !i.GetReturnOnlyChangedValues() || i.changedValues[k] { // change flag not set or object changed
			items[k] = v
		}
	}

	return items
}

func (i *InMemoryBackingStore) EnumerateKeysForValuesChangedToNil() []string {
	keys := make([]string, 0)
	for k, v := range i.store {
		valueOfV := reflect.ValueOf(v)
		if i.changedValues[k] && (v == nil || valueOfV.Kind() == reflect.Ptr && valueOfV.IsNil()) {
			keys = append(keys, k)
		}
	}

	return keys
}

func (i *InMemoryBackingStore) Subscribe(callback BackingStoreSubscriber) string {
	id := uuid.New().String()
	i.subscribers[id] = callback
	return id
}

func (i *InMemoryBackingStore) SubscribeWithId(callback BackingStoreSubscriber, subscriptionId string) error {
	subscriptionId = strings.TrimSpace(subscriptionId)
	if subscriptionId == "" {
		return errors.New("subscriptionId cannot be an empty string")
	}

	i.subscribers[subscriptionId] = callback

	return nil
}

func (i *InMemoryBackingStore) Unsubscribe(subscriptionId string) error {
	subscriptionId = strings.TrimSpace(subscriptionId)
	if subscriptionId == "" {
		return errors.New("subscriptionId cannot be an empty string")
	}

	delete(i.subscribers, subscriptionId)

	return nil
}

func (i *InMemoryBackingStore) Clear() {
	for k := range i.store {
		delete(i.store, k)
		delete(i.changedValues, k) // changed values must be an element in the store
	}
}

func (i *InMemoryBackingStore) GetInitializationCompleted() bool {
	return i.initializationCompleted
}

func (i *InMemoryBackingStore) SetInitializationCompleted(val bool) {
	i.initializationCompleted = val
}

func (i *InMemoryBackingStore) GetReturnOnlyChangedValues() bool {
	return i.returnOnlyChangedValues
}

func (i *InMemoryBackingStore) SetReturnOnlyChangedValues(val bool) {
	i.returnOnlyChangedValues = val
}
