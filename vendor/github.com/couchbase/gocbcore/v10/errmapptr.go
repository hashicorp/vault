package gocbcore

import (
	"sync/atomic"
	"unsafe"
)

type kvErrorMapPtr struct {
	data unsafe.Pointer
}

func (ptr *kvErrorMapPtr) Get() *kvErrorMap {
	return (*kvErrorMap)(atomic.LoadPointer(&ptr.data))
}

func (ptr *kvErrorMapPtr) Update(old, new *kvErrorMap) bool {
	if new == nil {
		logErrorf("Attempted to update to nil kvErrorMap")
		return false
	}

	return atomic.CompareAndSwapPointer(&ptr.data, unsafe.Pointer(old), unsafe.Pointer(new))
}
