// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/mitchellh/go-testing-interface"
	"github.com/mitchellh/mapstructure"
)

// ToMap renders an input value of any type as a map.  This is intended for
// logging human-readable data dumps in test logs, so it uses the `json`
// tags on struct fields: this makes it easy to exclude `"-"` values that
// are typically not interesting, respect omitempty, etc.
//
// We also replace any []byte fields with a hash of their value.
// This is usually sufficient for test log purposes, and is a lot more readable
// than a big array of individual byte values like Go would normally stringify a
// byte slice.
func ToMap(in any) (map[string]any, error) {
	temp := make(map[string]any)
	cfg := &mapstructure.DecoderConfig{
		TagName:              "json",
		IgnoreUntaggedFields: true,
		Result:               &temp,
	}
	md, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return nil, err
	}
	err = md.Decode(in)
	if err != nil {
		return nil, err
	}

	// mapstructure doesn't call the DecodeHook for each field when doing
	// struct->map conversions, but it does for map->map, so call it a second
	// time to convert each []byte field.
	out := make(map[string]any)
	md2, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &out,
		DecodeHook: func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
			if from.Kind() != reflect.Slice || from.Elem().Kind() != reflect.Uint8 {
				return data, nil
			}
			b := data.([]byte)
			return fmt.Sprintf("%x", sha256.Sum256(b)), nil
		},
	})
	if err != nil {
		return nil, err
	}
	err = md2.Decode(temp)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ToString renders its input using ToMap, and returns a string containing the
// result or an error if that fails.
func ToString(in any) string {
	m, err := ToMap(in)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", m)
}

// StringOrDie renders its input using ToMap, and returns a string containing the
// result.  If rendering yields an error, calls t.Fatal.
func StringOrDie(t testing.T, in any) string {
	t.Helper()
	m, err := ToMap(in)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%v", m)
}
