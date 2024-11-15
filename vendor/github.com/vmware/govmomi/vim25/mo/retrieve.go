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
	"context"
	"reflect"

	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

func ignoreMissingProperty(ref types.ManagedObjectReference, p types.MissingProperty) bool {
	switch ref.Type {
	case "VirtualMachine":
		switch p.Path {
		case "environmentBrowser":
			// See https://github.com/vmware/govmomi/pull/242
			return true
		case "alarmActionsEnabled":
			// Seen with vApp child VM
			return true
		}
	}

	return false
}

// ObjectContentToType loads an ObjectContent value into the value it
// represents. If the ObjectContent value has a non-empty 'MissingSet' field,
// it returns the first fault it finds there as error. If the 'MissingSet'
// field is empty, it returns a pointer to a reflect.Value. It handles contain
// nested properties, such as 'guest.ipAddress' or 'config.hardware'.
func ObjectContentToType(o types.ObjectContent, ptr ...bool) (interface{}, error) {
	// Expect no properties in the missing set
	for _, p := range o.MissingSet {
		if ignoreMissingProperty(o.Obj, p) {
			continue
		}

		return nil, soap.WrapVimFault(p.Fault.Fault)
	}

	ti := typeInfoForType(o.Obj.Type)
	v, err := ti.LoadFromObjectContent(o)
	if err != nil {
		return nil, err
	}

	if len(ptr) == 1 && ptr[0] {
		return v.Interface(), nil
	}
	return v.Elem().Interface(), nil
}

// ApplyPropertyChange converts the response of a call to WaitForUpdates
// and applies it to the given managed object.
func ApplyPropertyChange(obj Reference, changes []types.PropertyChange) {
	t := typeInfoForType(obj.Reference().Type)
	v := reflect.ValueOf(obj)

	for _, p := range changes {
		var field Field
		if !field.FromString(p.Name) {
			panic(p.Name + ": invalid property path")
		}

		rv, ok := t.props[field.Path]
		if !ok {
			panic(field.Path + ": property not found")
		}

		if field.Key == nil { // Key is only used for notifications
			assignValue(v, rv, reflect.ValueOf(p.Val))
		}
	}
}

// LoadObjectContent converts the response of a call to
// RetrieveProperties{Ex} to one or more managed objects.
func LoadObjectContent(content []types.ObjectContent, dst interface{}) error {
	rt := reflect.TypeOf(dst)
	if rt == nil || rt.Kind() != reflect.Ptr {
		panic("need pointer")
	}

	rv := reflect.ValueOf(dst).Elem()
	if !rv.CanSet() {
		panic("cannot set dst")
	}

	isSlice := false
	switch rt.Elem().Kind() {
	case reflect.Struct:
	case reflect.Slice:
		isSlice = true
	default:
		panic("unexpected type")
	}

	if isSlice {
		for _, p := range content {
			v, err := ObjectContentToType(p)
			if err != nil {
				return err
			}

			vt := reflect.TypeOf(v)

			if !rv.Type().AssignableTo(vt) {
				// For example: dst is []ManagedEntity, res is []HostSystem
				if field, ok := vt.FieldByName(rt.Elem().Elem().Name()); ok && field.Anonymous {
					rv.Set(reflect.Append(rv, reflect.ValueOf(v).FieldByIndex(field.Index)))
					continue
				}
			}

			rv.Set(reflect.Append(rv, reflect.ValueOf(v)))
		}
	} else {
		switch len(content) {
		case 0:
		case 1:
			v, err := ObjectContentToType(content[0])
			if err != nil {
				return err
			}

			vt := reflect.TypeOf(v)

			if !rv.Type().AssignableTo(vt) {
				// For example: dst is ComputeResource, res is ClusterComputeResource
				if field, ok := vt.FieldByName(rt.Elem().Name()); ok && field.Anonymous {
					rv.Set(reflect.ValueOf(v).FieldByIndex(field.Index))
					return nil
				}
			}

			rv.Set(reflect.ValueOf(v))
		default:
			// If dst is not a slice, expect to receive 0 or 1 results
			panic("more than 1 result")
		}
	}

	return nil
}

// RetrievePropertiesEx wraps RetrievePropertiesEx and ContinueRetrievePropertiesEx to collect properties in batches.
func RetrievePropertiesEx(ctx context.Context, r soap.RoundTripper, req types.RetrievePropertiesEx) ([]types.ObjectContent, error) {
	rx, err := methods.RetrievePropertiesEx(ctx, r, &req)
	if err != nil {
		return nil, err
	}

	if rx.Returnval == nil {
		return nil, nil
	}

	objects := rx.Returnval.Objects
	token := rx.Returnval.Token

	for token != "" {
		cx, err := methods.ContinueRetrievePropertiesEx(ctx, r, &types.ContinueRetrievePropertiesEx{
			This:  req.This,
			Token: token,
		})
		if err != nil {
			return nil, err
		}

		token = cx.Returnval.Token
		objects = append(objects, cx.Returnval.Objects...)
	}

	return objects, nil
}

// RetrievePropertiesForRequest calls the RetrieveProperties method with the
// specified request and decodes the response struct into the value pointed to
// by dst.
func RetrievePropertiesForRequest(ctx context.Context, r soap.RoundTripper, req types.RetrieveProperties, dst interface{}) error {
	objects, err := RetrievePropertiesEx(ctx, r, types.RetrievePropertiesEx{
		This:    req.This,
		SpecSet: req.SpecSet,
	})
	if err != nil {
		return err
	}

	return LoadObjectContent(objects, dst)
}

// RetrieveProperties retrieves the properties of the managed object specified
// as obj and decodes the response struct into the value pointed to by dst.
func RetrieveProperties(ctx context.Context, r soap.RoundTripper, pc, obj types.ManagedObjectReference, dst interface{}) error {
	req := types.RetrieveProperties{
		This: pc,
		SpecSet: []types.PropertyFilterSpec{
			{
				ObjectSet: []types.ObjectSpec{
					{
						Obj:  obj,
						Skip: types.NewBool(false),
					},
				},
				PropSet: []types.PropertySpec{
					{
						All:  types.NewBool(true),
						Type: obj.Type,
					},
				},
			},
		},
	}

	return RetrievePropertiesForRequest(ctx, r, req, dst)
}

var morType = reflect.TypeOf((*types.ManagedObjectReference)(nil)).Elem()

// References returns all non-nil moref field values in the given struct.
// Only Anonymous struct fields are followed by default. The optional follow
// param will follow any struct fields when true.
func References(s interface{}, follow ...bool) []types.ManagedObjectReference {
	var refs []types.ManagedObjectReference
	rval := reflect.ValueOf(s)
	rtype := rval.Type()

	if rval.Kind() == reflect.Ptr {
		rval = rval.Elem()
		rtype = rval.Type()
	}

	for i := 0; i < rval.NumField(); i++ {
		val := rval.Field(i)
		finfo := rtype.Field(i)

		if finfo.Anonymous {
			refs = append(refs, References(val.Interface(), follow...)...)
			continue
		}
		if finfo.Name == "Self" {
			continue
		}

		ftype := val.Type()

		if ftype.Kind() == reflect.Slice {
			if ftype.Elem() == morType {
				s := val.Interface().([]types.ManagedObjectReference)
				for i := range s {
					refs = append(refs, s[i])
				}
			}
			continue
		}

		if ftype.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
			ftype = val.Type()
		}

		if ftype == morType {
			refs = append(refs, val.Interface().(types.ManagedObjectReference))
			continue
		}

		if len(follow) != 0 && follow[0] {
			if ftype.Kind() == reflect.Struct && val.CanSet() {
				refs = append(refs, References(val.Interface(), follow...)...)
			}
		}
	}

	return refs
}
