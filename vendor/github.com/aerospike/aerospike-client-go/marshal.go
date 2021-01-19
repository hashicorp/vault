// +build !as_performance

// Copyright 2013-2019 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"
)

var aerospikeTag = "as"

const (
	aerospikeMetaTag    = "asm"
	aerospikeMetaTagGen = "gen"
	aerospikeMetaTagTTL = "ttl"
)

// This method is copied verbatim from https://golang.org/src/encoding/json/encode.go
// to ensure compatibility with the json package.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}

// SetAerospikeTag sets the bin tag to the specified tag.
// This will be useful for when a user wants to use the same tag name for two different concerns.
// For example, one will be able to use the same tag name for both json and aerospike bin name.
func SetAerospikeTag(tag string) {
	aerospikeTag = tag
}

func valueToInterface(f reflect.Value, clusterSupportsFloat bool) interface{} {
	// get to the core value
	for f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return nil
		}
		f = reflect.Indirect(f)
	}

	switch f.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return IntegerValue(f.Int())
	case reflect.Uint64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint16:
		return int64(f.Uint())
	case reflect.Float64, reflect.Float32:
		// support floats through integer encoding if
		// server doesn't support floats
		if clusterSupportsFloat {
			return FloatValue(f.Float())
		}
		return IntegerValue(math.Float64bits(f.Float()))

	case reflect.Struct:
		if f.Type().PkgPath() == "time" && f.Type().Name() == "Time" {
			return f.Interface().(time.Time).UTC().UnixNano()
		}
		return structToMap(f, clusterSupportsFloat)
	case reflect.Bool:
		if f.Bool() {
			return IntegerValue(1)
		}
		return IntegerValue(0)
	case reflect.Map:
		if f.IsNil() {
			return nil
		}

		newMap := make(map[interface{}]interface{}, f.Len())
		for _, mk := range f.MapKeys() {
			newMap[valueToInterface(mk, clusterSupportsFloat)] = valueToInterface(f.MapIndex(mk), clusterSupportsFloat)
		}

		return newMap
	case reflect.Slice, reflect.Array:
		if f.Kind() == reflect.Slice && f.IsNil() {
			return nil
		}
		if f.Kind() == reflect.Slice && reflect.TypeOf(f.Interface()).Elem().Kind() == reflect.Uint8 {
			// handle blobs
			return f.Interface().([]byte)
		}
		// convert to primitives recursively
		newSlice := make([]interface{}, f.Len(), f.Cap())
		for i := 0; i < len(newSlice); i++ {
			newSlice[i] = valueToInterface(f.Index(i), clusterSupportsFloat)
		}
		return newSlice
	case reflect.Interface:
		if f.IsNil() {
			return nullValue
		}
		return f.Interface()
	default:
		return f.Interface()
	}
}

func fieldIsMetadata(f reflect.StructField) bool {
	meta := f.Tag.Get(aerospikeMetaTag)
	return strings.Trim(meta, " ") != ""
}

func fieldIsOmitOnEmpty(f reflect.StructField) bool {
	tag := f.Tag.Get(aerospikeTag)
	return strings.Contains(tag, ",omitempty")
}

func stripOptions(tag string) string {
	i := strings.Index(tag, ",")
	if i < 0 {
		return tag
	}
	return string(tag[:i])
}

func fieldAlias(f reflect.StructField) string {
	alias := strings.Trim(stripOptions(f.Tag.Get(aerospikeTag)), " ")
	if alias != "" {
		// if tag is -, the field should not be persisted
		if alias == "-" {
			return ""
		}
		return alias
	}
	return f.Name
}

func setBinMap(s reflect.Value, clusterSupportsFloat bool, typeOfT reflect.Type, binMap BinMap, index []int) {
	numFields := typeOfT.NumField()
	var fld reflect.StructField
	for i := 0; i < numFields; i++ {
		fld = typeOfT.Field(i)

		fldIndex := append(index, fld.Index...)

		if fld.Anonymous && fld.Type.Kind() == reflect.Struct {
			setBinMap(s, clusterSupportsFloat, fld.Type, binMap, fldIndex)
			continue
		}

		// skip unexported fields
		if fld.PkgPath != "" {
			continue
		}

		if fieldIsMetadata(fld) {
			continue
		}

		// skip transient fields tagged `-`
		alias := fieldAlias(fld)
		if alias == "" {
			continue
		}

		value := s.FieldByIndex(fldIndex)
		if fieldIsOmitOnEmpty(fld) && isEmptyValue(value) {
			continue
		}

		binValue := valueToInterface(value, clusterSupportsFloat)

		if _, ok := binMap[alias]; ok {
			panic(fmt.Sprintf("ambiguous fields with the same name or alias: %s", alias))
		}
		binMap[alias] = binValue
	}
}

func structToMap(s reflect.Value, clusterSupportsFloat bool) BinMap {
	if !s.IsValid() {
		return nil
	}

	var binMap BinMap = make(BinMap, s.NumField())

	setBinMap(s, clusterSupportsFloat, s.Type(), binMap, nil)

	return binMap
}

func marshal(v interface{}, clusterSupportsFloat bool) BinMap {
	s := indirect(reflect.ValueOf(v))
	return structToMap(s, clusterSupportsFloat)
}

type syncMap struct {
	objectMappings map[reflect.Type]map[string][]int
	objectFields   map[reflect.Type][]string
	objectTTLs     map[reflect.Type][][]int
	objectGen      map[reflect.Type][][]int
	mutex          sync.RWMutex
}

func (sm *syncMap) setMapping(objType reflect.Type, mapping map[string][]int, fields []string, ttl, gen [][]int) {
	sm.mutex.Lock()
	sm.objectMappings[objType] = mapping
	sm.objectFields[objType] = fields
	sm.objectTTLs[objType] = ttl
	sm.objectGen[objType] = gen
	sm.mutex.Unlock()
}

func indirect(obj reflect.Value) reflect.Value {
	for obj.Kind() == reflect.Ptr {
		if obj.IsNil() {
			return obj
		}
		obj = obj.Elem()
	}
	return obj
}

func indirectT(objType reflect.Type) reflect.Type {
	for objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	return objType
}

func (sm *syncMap) mappingExists(objType reflect.Type) (map[string][]int, bool) {
	sm.mutex.RLock()
	mapping, exists := sm.objectMappings[objType]
	sm.mutex.RUnlock()
	return mapping, exists
}

func (sm *syncMap) getMapping(objType reflect.Type) map[string][]int {
	objType = indirectT(objType)
	mapping, exists := sm.mappingExists(objType)
	if !exists {
		cacheObjectTags(objType)
		mapping, _ = sm.mappingExists(objType)
	}

	return mapping
}

func (sm *syncMap) getMetaMappings(objType reflect.Type) (ttl, gen [][]int) {
	objType = indirectT(objType)
	if _, exists := sm.mappingExists(objType); !exists {
		cacheObjectTags(objType)
	}

	sm.mutex.RLock()
	ttl = sm.objectTTLs[objType]
	gen = sm.objectGen[objType]
	sm.mutex.RUnlock()
	return ttl, gen
}

func (sm *syncMap) fieldsExists(objType reflect.Type) ([]string, bool) {
	sm.mutex.RLock()
	mapping, exists := sm.objectFields[objType]
	sm.mutex.RUnlock()
	return mapping, exists
}

func (sm *syncMap) getFields(objType reflect.Type) []string {
	objType = indirectT(objType)
	fields, exists := sm.fieldsExists(objType)
	if !exists {
		cacheObjectTags(objType)
		fields, _ = sm.fieldsExists(objType)
	}

	return fields
}

var objectMappings = &syncMap{
	objectMappings: map[reflect.Type]map[string][]int{},
	objectFields:   map[reflect.Type][]string{},
	objectTTLs:     map[reflect.Type][][]int{},
	objectGen:      map[reflect.Type][][]int{},
}

func fillMapping(objType reflect.Type, mapping map[string][]int, fields []string, ttl, gen [][]int, index []int) ([]string, [][]int, [][]int) {
	numFields := objType.NumField()
	for i := 0; i < numFields; i++ {
		f := objType.Field(i)
		fIndex := append(index, f.Index...)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			fields, ttl, gen = fillMapping(f.Type, mapping, fields, ttl, gen, fIndex)
			continue
		}

		// skip unexported fields
		if f.PkgPath != "" {
			continue
		}

		tag := strings.Trim(stripOptions(f.Tag.Get(aerospikeTag)), " ")
		tagM := strings.Trim(f.Tag.Get(aerospikeMetaTag), " ")

		if tag != "" && tagM != "" {
			panic(fmt.Sprintf("Cannot accept both data and metadata tags on the same attribute on struct: %s.%s", objType.Name(), f.Name))
		}

		if tag != "-" && tagM == "" {
			if tag == "" {
				tag = f.Name
			}
			if _, ok := mapping[tag]; ok {
				panic(fmt.Sprintf("ambiguous fields with the same name or alias: %s", tag))
			}
			mapping[tag] = fIndex
			fields = append(fields, tag)
		}

		if tagM == aerospikeMetaTagTTL {
			ttl = append(ttl, fIndex)
		} else if tagM == aerospikeMetaTagGen {
			gen = append(gen, fIndex)
		} else if tagM != "" {
			panic(fmt.Sprintf("Invalid metadata tag `%s` on struct attribute: %s.%s", tagM, objType.Name(), f.Name))
		}
	}
	return fields, ttl, gen
}

func cacheObjectTags(objType reflect.Type) {
	mapping := map[string][]int{}
	fields, ttl, gen := fillMapping(objType, mapping, []string{}, nil, nil, nil)
	objectMappings.setMapping(objType, mapping, fields, ttl, gen)
}
