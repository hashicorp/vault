// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"fmt"
	"reflect"
)

type Filter struct {
	// The underlying boolean expression evaluator
	evaluator *Evaluator
}

// Creates a filter to operate on the given data type.
// The data type passed can be either be a container type (map, slice or array) or the element type.
// For example, if you want to filter a []Foo then the data type to pass here is either []Foo or just Foo.
// If no expression is provided the nil filter will be returned but is not an error. This is done
// to allow for executing the nil filter which is just a no-op
func CreateFilter(expression string) (*Filter, error) {
	if expression == "" {
		// nil filter
		return nil, nil
	}
	exp, err := CreateEvaluator(expression)
	if err != nil {
		return nil, fmt.Errorf("Failed to create boolean expression evaluator: %v", err)
	}

	return &Filter{
		evaluator: exp,
	}, nil
}

// Execute the filter. If called on a nil filter this is a no-op and
// will return the original data
func (f *Filter) Execute(data interface{}) (interface{}, error) {
	if f == nil {
		return data, nil
	}

	rvalue := reflect.ValueOf(data)
	rtype := rvalue.Type()

	switch rvalue.Kind() {
	case reflect.Array:
		// For arrays we return slices instead of fixed sized arrays
		rtype = reflect.SliceOf(rtype.Elem())
		fallthrough
	case reflect.Slice:
		newSlice := reflect.MakeSlice(rtype, 0, rvalue.Len())

		for i := 0; i < rvalue.Len(); i++ {
			item := rvalue.Index(i)
			if !item.CanInterface() {
				return nil, fmt.Errorf("Slice/Array value can not be used")
			}
			result, err := f.evaluator.Evaluate(item.Interface())
			if err != nil {
				return nil, err
			}

			if result {
				newSlice = reflect.Append(newSlice, item)
			}
		}

		return newSlice.Interface(), nil
	case reflect.Map:
		newMap := reflect.MakeMap(rtype)

		// TODO (mkeeler) - Update to use a MapRange iterator once Go 1.12 is usable
		// for all of our products
		for _, mapKey := range rvalue.MapKeys() {
			item := rvalue.MapIndex(mapKey)

			if !item.CanInterface() {
				return nil, fmt.Errorf("Map value cannot be used")
			}

			result, err := f.evaluator.Evaluate(item.Interface())
			if err != nil {
				return nil, err
			}

			if result {
				newMap.SetMapIndex(mapKey, item)
			}
		}

		return newMap.Interface(), nil
	default:
		return nil, fmt.Errorf("Only slices, arrays and maps are filterable")
	}
}
