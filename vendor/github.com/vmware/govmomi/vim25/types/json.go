/*
Copyright (c) 2023-2023 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"bytes"
	"io"
	"reflect"
	"time"

	"github.com/vmware/govmomi/vim25/json"
)

const (
	discriminatorMemberName  = "_typeName"
	primitiveValueMemberName = "_value"
)

var discriminatorTypeRegistry = map[string]reflect.Type{
	"boolean":  reflect.TypeOf(true),
	"byte":     reflect.TypeOf(uint8(0)),
	"short":    reflect.TypeOf(int16(0)),
	"int":      reflect.TypeOf(int32(0)),
	"long":     reflect.TypeOf(int64(0)),
	"float":    reflect.TypeOf(float32(0)),
	"double":   reflect.TypeOf(float64(0)),
	"string":   reflect.TypeOf(""),
	"binary":   reflect.TypeOf([]byte{}),
	"dateTime": reflect.TypeOf(time.Now()),
}

// NewJSONDecoder creates JSON decoder configured for VMOMI.
func NewJSONDecoder(r io.Reader) *json.Decoder {
	res := json.NewDecoder(r)
	res.SetDiscriminator(
		discriminatorMemberName,
		primitiveValueMemberName,
		json.DiscriminatorToTypeFunc(func(name string) (reflect.Type, bool) {
			if res, ok := TypeFunc()(name); ok {
				return res, true
			}
			if res, ok := discriminatorTypeRegistry[name]; ok {
				return res, true
			}
			return nil, false
		}),
	)
	return res
}

// VMOMI primitive names
var discriminatorNamesRegistry = map[reflect.Type]string{
	reflect.TypeOf(true):       "boolean",
	reflect.TypeOf(uint8(0)):   "byte",
	reflect.TypeOf(int16(0)):   "short",
	reflect.TypeOf(int32(0)):   "int",
	reflect.TypeOf(int64(0)):   "long",
	reflect.TypeOf(float32(0)): "float",
	reflect.TypeOf(float64(0)): "double",
	reflect.TypeOf(""):         "string",
	reflect.TypeOf([]byte{}):   "binary",
	reflect.TypeOf(time.Now()): "dateTime",
}

// NewJSONEncoder creates JSON encoder configured for VMOMI.
func NewJSONEncoder(w *bytes.Buffer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetDiscriminator(
		discriminatorMemberName,
		primitiveValueMemberName,
		json.DiscriminatorEncodeTypeNameRootValue|
			json.DiscriminatorEncodeTypeNameAllObjects,
	)
	enc.SetTypeToDiscriminatorFunc(VmomiTypeName)
	return enc
}

// VmomiTypeName computes the VMOMI type name of a go type. It uses a lookup
// table for VMOMI primitive types and the default discriminator function for
// other types.
func VmomiTypeName(t reflect.Type) (discriminator string) {
	// Look up primitive type names from VMOMI protocol
	if name, ok := discriminatorNamesRegistry[t]; ok {
		return name
	}
	name := json.DefaultDiscriminatorFunc(t)
	return name
}
