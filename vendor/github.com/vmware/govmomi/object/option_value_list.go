/*
Copyright (c) 2024-2024 VMware, Inc. All Rights Reserved.

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

package object

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/vmware/govmomi/vim25/types"
)

// OptionValueList simplifies manipulation of properties that are arrays of
// types.BaseOptionValue, such as ExtraConfig.
type OptionValueList []types.BaseOptionValue

// OptionValueListFromMap returns a new OptionValueList object from the provided
// map.
func OptionValueListFromMap[T any](in map[string]T) OptionValueList {
	if len(in) == 0 {
		return nil
	}
	var (
		i   int
		out = make(OptionValueList, len(in))
	)
	for k, v := range in {
		out[i] = &types.OptionValue{Key: k, Value: v}
		i++
	}
	return out
}

// IsTrue returns true if the specified key exists with an empty value or value
// equal to 1, "1", "on", "t", true, "true", "y", or "yes".
// All string comparisons are case-insensitive.
func (ov OptionValueList) IsTrue(key string) bool {
	return ov.isTrueOrFalse(key, true, 1, "", "1", "on", "t", "true", "y", "yes")
}

// IsFalse returns true if the specified key exists and has a value equal to
// 0, "0", "f", false, "false", "n", "no", or "off".
// All string comparisons are case-insensitive.
func (ov OptionValueList) IsFalse(key string) bool {
	return ov.isTrueOrFalse(key, false, 0, "0", "f", "false", "n", "no", "off")
}

func (ov OptionValueList) isTrueOrFalse(
	key string,
	boolVal bool,
	numVal int,
	strVals ...string) bool {

	val, ok := ov.Get(key)
	if !ok {
		return false
	}

	switch tval := val.(type) {
	case string:
		return slices.Contains(strVals, strings.ToLower(tval))
	case bool:
		return tval == boolVal
	case uint:
		return tval == uint(numVal)
	case uint8:
		return tval == uint8(numVal)
	case uint16:
		return tval == uint16(numVal)
	case uint32:
		return tval == uint32(numVal)
	case uint64:
		return tval == uint64(numVal)
	case int:
		return tval == int(numVal)
	case int8:
		return tval == int8(numVal)
	case int16:
		return tval == int16(numVal)
	case int32:
		return tval == int32(numVal)
	case int64:
		return tval == int64(numVal)
	case float32:
		return tval == float32(numVal)
	case float64:
		return tval == float64(numVal)
	}

	return false
}

// Get returns the value if exists, otherwise nil is returned. The second return
// value is a flag indicating whether the value exists or nil was the actual
// value.
func (ov OptionValueList) Get(key string) (any, bool) {
	if ov == nil {
		return nil, false
	}
	for i := range ov {
		if optVal := ov[i].GetOptionValue(); optVal != nil {
			if optVal.Key == key {
				return optVal.Value, true
			}
		}
	}
	return nil, false
}

// GetString returns the value as a string if the value exists.
func (ov OptionValueList) GetString(key string) (string, bool) {
	if ov == nil {
		return "", false
	}
	for i := range ov {
		if optVal := ov[i].GetOptionValue(); optVal != nil {
			if optVal.Key == key {
				return getOptionValueAsString(optVal.Value), true
			}
		}
	}
	return "", false
}

// Additions returns a diff that includes only the elements from the provided
// list that do not already exist.
func (ov OptionValueList) Additions(in ...types.BaseOptionValue) OptionValueList {
	return ov.diff(in, true)
}

// Diff returns a diff that includes the elements from the provided list that do
// not already exist or have different values.
func (ov OptionValueList) Diff(in ...types.BaseOptionValue) OptionValueList {
	return ov.diff(in, false)
}

func (ov OptionValueList) diff(in OptionValueList, addOnly bool) OptionValueList {
	if ov == nil && in == nil {
		return nil
	}
	var (
		out         OptionValueList
		leftOptVals = ov.Map()
	)
	for i := range in {
		if rightOptVal := in[i].GetOptionValue(); rightOptVal != nil {
			k, v := rightOptVal.Key, rightOptVal.Value
			if ov == nil {
				out = append(out, &types.OptionValue{Key: k, Value: v})
			} else if leftOptVal, ok := leftOptVals[k]; !ok {
				out = append(out, &types.OptionValue{Key: k, Value: v})
			} else if !addOnly && v != leftOptVal {
				out = append(out, &types.OptionValue{Key: k, Value: v})
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// Join combines this list with the provided one and returns the result, joining
// the two lists on their shared keys.
// Please note, Join(left, right) means the values from right will be appended
// to left, without overwriting any values that have shared keys. To overwrite
// the shared keys in left from right, use Join(right, left) instead.
func (ov OptionValueList) Join(in ...types.BaseOptionValue) OptionValueList {
	var (
		out     OptionValueList
		outKeys map[string]struct{}
	)

	// Init the out slice from the left side.
	if len(ov) > 0 {
		outKeys = map[string]struct{}{}
		for i := range ov {
			if optVal := ov[i].GetOptionValue(); optVal != nil {
				kv := &types.OptionValue{Key: optVal.Key, Value: optVal.Value}
				out = append(out, kv)
				outKeys[optVal.Key] = struct{}{}
			}
		}
	}

	// Join the values from the right side.
	for i := range in {
		if rightOptVal := in[i].GetOptionValue(); rightOptVal != nil {
			k, v := rightOptVal.Key, rightOptVal.Value
			if _, ok := outKeys[k]; !ok {
				out = append(out, &types.OptionValue{Key: k, Value: v})
			}
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}

// Map returns the list of option values as a map. A nil value is returned if
// the list is empty.
func (ov OptionValueList) Map() map[string]any {
	if len(ov) == 0 {
		return nil
	}
	out := map[string]any{}
	for i := range ov {
		if optVal := ov[i].GetOptionValue(); optVal != nil {
			out[optVal.Key] = optVal.Value
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// StringMap returns the list of option values as a map where the values are
// strings. A nil value is returned if the list is empty.
func (ov OptionValueList) StringMap() map[string]string {
	if len(ov) == 0 {
		return nil
	}
	out := map[string]string{}
	for i := range ov {
		if optVal := ov[i].GetOptionValue(); optVal != nil {
			out[optVal.Key] = getOptionValueAsString(optVal.Value)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func getOptionValueAsString(val any) string {
	switch tval := val.(type) {
	case string:
		return tval
	default:
		if rv := reflect.ValueOf(val); rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return ""
			}
			return fmt.Sprintf("%v", rv.Elem().Interface())
		}
		return fmt.Sprintf("%v", tval)
	}
}
