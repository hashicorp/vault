// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package array

import (
	"reflect"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/memory"
	"golang.org/x/xerrors"
)

// ExtensionArray is the interface that needs to be implemented to handle
// user-defined extension type arrays. In order to ensure consistency and
// proper behavior, all ExtensionArray types must embed ExtensionArrayBase
// in order to meet the interface which provides the default implementation
// and handling for the array while allowing custom behavior to be built
// on top of it.
type ExtensionArray interface {
	Interface
	// ExtensionType returns the datatype as per calling DataType(), but
	// already cast to ExtensionType
	ExtensionType() arrow.ExtensionType
	// Storage returns the underlying storage array for this array.
	Storage() Interface

	// by having a non-exported function in the interface, it means that
	// consumers must embed ExtensionArrayBase in their structs in order
	// to fulfill this interface.
	mustEmbedExtensionArrayBase()
}

// two extension arrays are equal if their data types are equal and
// their underlying storage arrays are equal.
func arrayEqualExtension(l, r ExtensionArray) bool {
	if !arrow.TypeEqual(l.DataType(), r.DataType()) {
		return false
	}

	return ArrayEqual(l.Storage(), r.Storage())
}

// two extension arrays are approximately equal if their data types are
// equal and their underlying storage arrays are approximately equal.
func arrayApproxEqualExtension(l, r ExtensionArray, opt equalOption) bool {
	if !arrow.TypeEqual(l.DataType(), r.DataType()) {
		return false
	}

	return arrayApproxEqual(l.Storage(), r.Storage(), opt)
}

// NewExtensionArrayWithStorage constructs a new ExtensionArray from the provided
// ExtensionType and uses the provided storage interface as the underlying storage.
// This will not release the storage array passed in so consumers should call Release
// on it manually while the new Extension array will share references to the underlying
// Data buffers.
func NewExtensionArrayWithStorage(dt arrow.ExtensionType, storage Interface) Interface {
	if !arrow.TypeEqual(dt.StorageType(), storage.DataType()) {
		panic(xerrors.Errorf("arrow/array: storage type %s for extension type %s, does not match expected type %s", storage.DataType(), dt.ExtensionName(), dt.StorageType()))
	}

	base := ExtensionArrayBase{}
	base.refCount = 1
	base.storage = storage
	storage.Retain()

	storageData := storage.Data()
	// create a new data instance with the ExtensionType as the datatype but referencing the
	// same underlying buffers to share them with the storage array.
	baseData := NewData(dt, storageData.length, storageData.buffers, storageData.childData, storageData.nulls, storageData.offset)
	defer baseData.Release()
	base.array.setData(baseData)

	// use the ExtensionType's ArrayType to construct the correctly typed object
	// to use as the ExtensionArray interface. reflect.New returns a pointer to
	// the newly created object.
	arr := reflect.New(base.ExtensionType().ArrayType())
	// set the embedded ExtensionArrayBase to the value we created above. We know
	// that this field will exist because the interface requires embedding ExtensionArrayBase
	// so we don't have to separately check, this will panic if called on an ArrayType
	// that doesn't embed ExtensionArrayBase which is what we want.
	arr.Elem().FieldByName("ExtensionArrayBase").Set(reflect.ValueOf(base))
	return arr.Interface().(ExtensionArray)
}

// NewExtensionData expects a data with a datatype of arrow.ExtensionType and
// underlying data built for the storage array.
func NewExtensionData(data *Data) ExtensionArray {
	base := ExtensionArrayBase{}
	base.refCount = 1
	base.setData(data)

	// use the ExtensionType's ArrayType to construct the correctly typed object
	// to use as the ExtensionArray interface. reflect.New returns a pointer to
	// the newly created object.
	arr := reflect.New(base.ExtensionType().ArrayType())
	// set the embedded ExtensionArrayBase to the value we created above. We know
	// that this field will exist because the interface requires embedding ExtensionArrayBase
	// so we don't have to separately check, this will panic if called on an ArrayType
	// that doesn't embed ExtensionArrayBase which is what we want.
	arr.Elem().FieldByName("ExtensionArrayBase").Set(reflect.ValueOf(base))
	return arr.Interface().(ExtensionArray)
}

// ExtensionArrayBase is the base struct for user-defined Extension Array types
// and must be embedded in any user-defined extension arrays like so:
//
//   type UserDefinedArray struct {
//       array.ExtensionArrayBase
//   }
//
type ExtensionArrayBase struct {
	array
	storage Interface
}

// Retain increases the reference count by 1.
// Retain may be called simultaneously from multiple goroutines.
func (e *ExtensionArrayBase) Retain() {
	e.array.Retain()
	e.storage.Retain()
}

// Release decreases the reference count by 1.
// Release may be called simultaneously from multiple goroutines.
// When the reference count goes to zero, the memory is freed.
func (e *ExtensionArrayBase) Release() {
	e.array.Release()
	e.storage.Release()
}

// Storage returns the underlying storage array
func (e *ExtensionArrayBase) Storage() Interface { return e.storage }

// ExtensionType returns the same thing as DataType, just already casted
// to an ExtensionType interface for convenience.
func (e *ExtensionArrayBase) ExtensionType() arrow.ExtensionType {
	return e.DataType().(arrow.ExtensionType)
}

func (e *ExtensionArrayBase) setData(data *Data) {
	if data.DataType().ID() != arrow.EXTENSION {
		panic("arrow/array: must use extension type to construct an extension array")
	}
	extType, ok := data.dtype.(arrow.ExtensionType)
	if !ok {
		panic("arrow/array: DataType for ExtensionArray must implement arrow.ExtensionType")
	}

	e.array.setData(data)
	// our underlying storage needs to reference the same data buffers (no copying)
	// but should have the storage type's datatype, so we create a Data for it.
	storageData := NewData(extType.StorageType(), data.length, data.buffers, data.childData, data.nulls, data.offset)
	defer storageData.Release()
	e.storage = MakeFromData(storageData)
}

// no-op function that exists simply to force embedding this in any extension array types.
func (ExtensionArrayBase) mustEmbedExtensionArrayBase() {}

// ExtensionBuilder is a convenience builder so that NewBuilder and such will still work
// with extension types properly. Depending on preference it may be cleaner or easier to just use
// NewExtensionArrayWithStorage and pass a storage array.
//
// That said, this allows easily building an extension array by providing the extension
// type and retrieving the storage builder.
type ExtensionBuilder struct {
	Builder
	dt arrow.ExtensionType
}

// NewExtensionBuilder returns a builder using the provided memory allocator for the desired
// extension type. It will internally construct a builder of the storage type for the extension
// type and keep a copy of the extension type. The underlying type builder can then be retrieved
// by calling `StorageBuilder` on this and then type asserting it to the desired builder type.
//
// After using the storage builder, calling NewArray or NewExtensionArray will construct
// the appropriate extension array type and set the storage correctly, resetting the builder for
// reuse.
//
// Example
//
// Simple example assuming an extension type of a UUID defined as a FixedSizeBinary(16) was registered
// using the type name "uuid":
//
//   uuidType := arrow.GetExtensionType("uuid")
//   bldr := array.NewExtensionBuilder(memory.DefaultAllocator, uuidType)
//   defer bldr.Release()
//   uuidBldr := bldr.StorageBuilder().(*array.FixedSizeBinaryBuilder)
//   /* build up the fixed size binary array as usual via Append/AppendValues */
//   uuidArr := bldr.NewExtensionArray()
//   defer uuidArr.Release()
//
// Because the storage builder is embedded in the Extension builder it also means
// that any of the functions available on the Builder interface can be called on
// an instance of ExtensionBuilder and will respond appropriately as the storage
// builder would for generically grabbing the Lenth, Cap, Nulls, reserving, etc.
func NewExtensionBuilder(mem memory.Allocator, dt arrow.ExtensionType) *ExtensionBuilder {
	return &ExtensionBuilder{Builder: NewBuilder(mem, dt.StorageType()), dt: dt}
}

// StorageBuilder returns the builder for the underlying storage type.
func (b *ExtensionBuilder) StorageBuilder() Builder { return b.Builder }

// NewArray creates a new array from the memory buffers used by the builder
// and resets the builder so it can be used to build a new array.
func (b *ExtensionBuilder) NewArray() Interface {
	return b.NewExtensionArray()
}

// NewExtensionArray creates an Extension array from the memory buffers used
// by the builder and resets the ExtensionBuilder so it can be used to build
// a new ExtensionArray of the same type.
func (b *ExtensionBuilder) NewExtensionArray() ExtensionArray {
	storage := b.Builder.NewArray()
	defer storage.Release()

	data := NewData(b.dt, storage.Len(), storage.Data().buffers, storage.Data().childData, storage.Data().nulls, 0)
	defer data.Release()
	return NewExtensionData(data)
}
