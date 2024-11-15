// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"reflect"
	"strconv"
)

// Each returns a validation rule that loops through an iterable (map, slice or array)
// and validates each value inside with the provided rules.
// An empty iterable is considered valid. Use the Required rule to make sure the iterable is not empty.
func Each(rules ...Rule) *EachRule {
	return &EachRule{
		rules: rules,
	}
}

type EachRule struct {
	rules []Rule
}

// Loops through the given iterable and calls the Ozzo Validate() method for each value.
func (r *EachRule) Validate(value interface{}) error {
	errs := Errors{}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			val := r.getInterface(v.MapIndex(k))
			if err := Validate(val, r.rules...); err != nil {
				errs[r.getString(k)] = err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			val := r.getInterface(v.Index(i))
			if err := Validate(val, r.rules...); err != nil {
				errs[strconv.Itoa(i)] = err
			}
		}
	default:
		return errors.New("must be an iterable (map, slice or array)")
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r *EachRule) getInterface(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	default:
		return value.Interface()
	}
}

func (r *EachRule) getString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return ""
		}
		return value.Elem().String()
	default:
		return value.String()
	}
}
