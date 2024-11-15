/*
Copyright (c) 2014-2024 VMware, Inc. All Rights Reserved.

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

package mo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/vmware/govmomi/vim25/types"
)

type typeInfo struct {
	typ reflect.Type

	// Field indices of "Self" field.
	self []int

	// Map property names to field indices.
	props map[string][]int
}

var typeInfoLock sync.RWMutex
var typeInfoMap = make(map[string]*typeInfo)

func typeInfoForType(tname string) *typeInfo {
	typeInfoLock.RLock()
	ti, ok := typeInfoMap[tname]
	typeInfoLock.RUnlock()

	if ok {
		return ti
	}

	// Create new typeInfo for type.
	if typ, ok := t[tname]; !ok {
		panic("unknown type: " + tname)
	} else {
		// Multiple routines may race to set it, but the result is the same.
		typeInfoLock.Lock()
		ti = newTypeInfo(typ)
		typeInfoMap[tname] = ti
		typeInfoLock.Unlock()
	}

	return ti
}

func baseType(ftyp reflect.Type) reflect.Type {
	base := strings.TrimPrefix(ftyp.Name(), "Base")
	switch base {
	case "MethodFault":
		return nil
	}
	if kind, ok := types.TypeFunc()(base); ok {
		return kind
	}
	return nil
}

func newTypeInfo(typ reflect.Type) *typeInfo {
	t := typeInfo{
		typ:   typ,
		props: make(map[string][]int),
	}

	t.build(typ, "", []int{})

	return &t
}

var managedObjectRefType = reflect.TypeOf((*types.ManagedObjectReference)(nil)).Elem()

func buildName(fn string, f reflect.StructField) string {
	if fn != "" {
		fn += "."
	}

	motag := f.Tag.Get("json")
	if motag != "" {
		tokens := strings.Split(motag, ",")
		if tokens[0] != "" {
			return fn + tokens[0]
		}
	}

	xmltag := f.Tag.Get("xml")
	if xmltag != "" {
		tokens := strings.Split(xmltag, ",")
		if tokens[0] != "" {
			return fn + tokens[0]
		}
	}

	return ""
}

func (t *typeInfo) build(typ reflect.Type, fn string, fi []int) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic("need struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		ftyp := f.Type

		// Copy field indices so they can be passed along.
		fic := make([]int, len(fi)+1)
		copy(fic, fi)
		fic[len(fi)] = i

		// Recurse into embedded field.
		if f.Anonymous {
			t.build(ftyp, fn, fic)
			continue
		}

		// Top level type has a "Self" field.
		if f.Name == "Self" && ftyp == managedObjectRefType {
			t.self = fic
			continue
		}

		fnc := buildName(fn, f)
		if fnc == "" {
			continue
		}

		t.props[fnc] = fic

		// Dereference pointer.
		if ftyp.Kind() == reflect.Ptr {
			ftyp = ftyp.Elem()
		}

		// Slices are not addressable by `foo.bar.qux`.
		if ftyp.Kind() == reflect.Slice {
			continue
		}

		// Skip the managed reference type.
		if ftyp == managedObjectRefType {
			continue
		}

		// Recurse into structs.
		if ftyp.Kind() == reflect.Struct {
			t.build(ftyp, fnc, fic)
		}

		// Base type can only access base fields, for example Datastore.Info
		// is types.BaseDataStore, so we create a new(types.DatastoreInfo)
		// Indexed property path may traverse into array element fields.
		// When interface, use the base type to index fields.
		// For example, BaseVirtualDevice:
		//   config.hardware.device[4000].deviceInfo.label
		if ftyp.Kind() == reflect.Interface {
			if base := baseType(ftyp); base != nil {
				t.build(base, fnc, fic)
			}
		}
	}
}

var nilValue reflect.Value

// assignValue assigns a value 'pv' to the struct pointed to by 'val', given a
// slice of field indices. It recurses into the struct until it finds the field
// specified by the indices. It creates new values for pointer types where
// needed.
func assignValue(val reflect.Value, fi []int, pv reflect.Value, field ...string) {
	// Indexed property path can only use base types
	if val.Kind() == reflect.Interface {
		base := baseType(val.Type())
		val.Set(reflect.New(base))
		val = val.Elem()
	}

	// Create new value if necessary.
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}

		val = val.Elem()
	}

	rv := val.Field(fi[0])
	fi = fi[1:]
	if len(fi) == 0 {
		if pv == nilValue {
			pv = reflect.Zero(rv.Type())
			rv.Set(pv)
			return
		}
		rt := rv.Type()
		pt := pv.Type()

		// If type is a pointer, create new instance of type.
		if rt.Kind() == reflect.Ptr {
			rv.Set(reflect.New(rt.Elem()))
			rv = rv.Elem()
			rt = rv.Type()
		}

		// If the target type is a slice, but the source is not, deference any ArrayOfXYZ type
		if rt.Kind() == reflect.Slice && pt.Kind() != reflect.Slice {
			if pt.Kind() == reflect.Ptr {
				pv = pv.Elem()
				pt = pt.Elem()
			}

			m := arrayOfRegexp.FindStringSubmatch(pt.Name())
			if len(m) > 0 {
				pv = pv.FieldByName(m[1]) // ArrayOfXYZ type has single field named XYZ
				pt = pv.Type()

				if !pv.IsValid() {
					panic(fmt.Sprintf("expected %s type to have field %s", m[0], m[1]))
				}
			}
		}

		// If type is an interface, check if pv implements it.
		if rt.Kind() == reflect.Interface && !pt.Implements(rt) {
			// Check if pointer to pv implements it.
			if reflect.PtrTo(pt).Implements(rt) {
				npv := reflect.New(pt)
				npv.Elem().Set(pv)
				pv = npv
				pt = pv.Type()
			} else {
				panic(fmt.Sprintf("type %s doesn't implement %s", pt.Name(), rt.Name()))
			}
		} else if rt.Kind() == reflect.Struct && pt.Kind() == reflect.Ptr {
			pv = pv.Elem()
			pt = pv.Type()
		}

		if pt.AssignableTo(rt) {
			rv.Set(pv)
		} else if rt.ConvertibleTo(pt) {
			rv.Set(pv.Convert(rt))
		} else if rt.Kind() == reflect.Slice {
			// Indexed array value
			path := field[0]
			isInterface := rt.Elem().Kind() == reflect.Interface

			if len(path) == 0 {
				// Append item (pv) directly to the array, converting to pointer if interface
				if isInterface {
					npv := reflect.New(pt)
					npv.Elem().Set(pv)
					pv = npv
					pt = pv.Type()
				}
			} else {
				// Construct item to be appended to the array, setting field within to value of pv
				var item reflect.Value
				if isInterface {
					base := baseType(rt.Elem())
					item = reflect.New(base)
				} else {
					item = reflect.New(rt.Elem())
				}

				field := newTypeInfo(item.Type())
				if ix, ok := field.props[path]; ok {
					assignValue(item, ix, pv)
				}

				if rt.Elem().Kind() == reflect.Struct {
					pv = item.Elem()
				} else {
					pv = item
				}
				pt = pv.Type()
			}

			rv.Set(reflect.Append(rv, pv))
		} else {
			panic(fmt.Sprintf("cannot assign %q (%s) to %q (%s)", rt.Name(), rt.Kind(), pt.Name(), pt.Kind()))
		}

		return
	}

	assignValue(rv, fi, pv, field...)
}

var arrayOfRegexp = regexp.MustCompile("ArrayOf(.*)$")

// LoadObjectFromContent loads properties from the 'PropSet' field in the
// specified ObjectContent value into the value it represents, which is
// returned as a reflect.Value.
func (t *typeInfo) LoadFromObjectContent(o types.ObjectContent) (reflect.Value, error) {
	v := reflect.New(t.typ)
	assignValue(v, t.self, reflect.ValueOf(o.Obj))

	for _, p := range o.PropSet {
		var field Field
		field.FromString(p.Name)

		rv, ok := t.props[field.Path]
		if !ok {
			continue
		}
		assignValue(v, rv, reflect.ValueOf(p.Val), field.Item)
	}

	return v, nil
}

func IsManagedObjectType(kind string) bool {
	_, ok := t[kind]
	return ok
}

// Value returns a new mo instance of the given ref Type.
func Value(ref types.ManagedObjectReference) (Reference, bool) {
	if rt, ok := t[ref.Type]; ok {
		val := reflect.New(rt)
		val.Interface().(Entity).Entity().Self = ref
		return val.Elem().Interface().(Reference), true
	}
	return nil, false
}

// Field of a ManagedObject in string form.
type Field struct {
	Path string
	Key  any
	Item string
}

func (f *Field) String() string {
	if f.Key == nil {
		return f.Path
	}

	var key, item string

	switch f.Key.(type) {
	case string:
		key = fmt.Sprintf("%q", f.Key)
	default:
		key = fmt.Sprintf("%d", f.Key)
	}

	if f.Item != "" {
		item = "." + f.Item
	}

	return fmt.Sprintf("%s[%s]%s", f.Path, key, item)
}

func (f *Field) FromString(spec string) bool {
	s := strings.SplitN(spec, "[", 2)
	f.Path = s[0]
	f.Key = nil
	f.Item = ""
	if len(s) == 1 {
		return true
	}

	parts := strings.SplitN(s[1], "]", 2)

	if len(parts) != 2 {
		return false
	}

	ix := strings.Trim(parts[0], `"`)

	if ix == parts[0] {
		v, err := strconv.ParseInt(ix, 0, 32)
		if err != nil {
			return false
		}
		f.Key = int32(v)
	} else {
		f.Key = ix
	}

	if parts[1] == "" {
		return true
	}

	if parts[1][0] != '.' {
		return false
	}
	f.Item = parts[1][1:]

	return true
}
