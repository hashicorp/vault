package atomic

import "sync"

// SyncVal allows synchronized access to a value
type SyncVal struct {
	val  interface{}
	lock sync.RWMutex
}

// NewSyncVal creates a new instance of SyncVal
func NewSyncVal(val interface{}) *SyncVal {
	return &SyncVal{val: val}
}

// Set updates the value of SyncVal with the passed argument
func (sv *SyncVal) Set(val interface{}) {
	sv.lock.Lock()
	sv.val = val
	sv.lock.Unlock()
}

// Get returns the value inside the SyncVal
func (sv *SyncVal) Get() interface{} {
	sv.lock.RLock()
	val := sv.val
	sv.lock.RUnlock()
	return val
}

// GetSyncedVia returns the value returned by the function f.
func (sv *SyncVal) GetSyncedVia(f func(interface{}) (interface{}, error)) (interface{}, error) {
	sv.lock.RLock()
	defer sv.lock.RUnlock()

	val, err := f(sv.val)
	return val, err
}

// Update gets a function and passes the value of SyncVal to it.
// If the resulting err is nil, it will update the value of SyncVal.
// It will return the resulting error to the caller.
func (sv *SyncVal) Update(f func(interface{}) (interface{}, error)) error {
	sv.lock.Lock()
	defer sv.lock.Unlock()

	val, err := f(sv.val)
	if err == nil {
		sv.val = val
	}
	return err
}
