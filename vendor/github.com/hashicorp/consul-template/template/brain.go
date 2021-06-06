package template

import (
	"sync"

	dep "github.com/hashicorp/consul-template/dependency"
)

// Brain is what Template uses to determine the values that are
// available for template parsing.
type Brain struct {
	sync.RWMutex

	// data is the map of individual dependencies and the most recent data for
	// that dependency.
	data map[string]interface{}

	// receivedData is an internal tracker of which dependencies have stored data
	// in the brain.
	receivedData map[string]struct{}
}

// NewBrain creates a new Brain with empty values for each
// of the key structs.
func NewBrain() *Brain {
	return &Brain{
		data:         make(map[string]interface{}),
		receivedData: make(map[string]struct{}),
	}
}

// Remember accepts a dependency and the data to store associated with that
// dep. This function converts the given data to a proper type and stores
// it interally.
func (b *Brain) Remember(d dep.Dependency, data interface{}) {
	b.Lock()
	defer b.Unlock()

	b.data[d.String()] = data
	b.receivedData[d.String()] = struct{}{}
}

// Recall gets the current value for the given dependency in the Brain.
func (b *Brain) Recall(d dep.Dependency) (interface{}, bool) {
	b.RLock()
	defer b.RUnlock()

	// If we have not received data for this dependency, return now.
	if _, ok := b.receivedData[d.String()]; !ok {
		return nil, false
	}

	return b.data[d.String()], true
}

// ForceSet is used to force set the value of a dependency
// for a given hash code
func (b *Brain) ForceSet(hashCode string, data interface{}) {
	b.Lock()
	defer b.Unlock()

	b.data[hashCode] = data
	b.receivedData[hashCode] = struct{}{}
}

// Forget accepts a dependency and removes all associated data with this
// dependency. It also resets the "receivedData" internal map.
func (b *Brain) Forget(d dep.Dependency) {
	b.Lock()
	defer b.Unlock()

	delete(b.data, d.String())
	delete(b.receivedData, d.String())
}
