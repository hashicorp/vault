/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

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
	"reflect"
	"strings"
)

var (
	t = map[string]reflect.Type{}

	// minAPIVersionForType is used to lookup the minimum API version for which
	// a type is valid.
	minAPIVersionForType = map[string]string{}

	// minAPIVersionForEnumValue is used to lookup the minimum API version for
	// which an enum value is valid.
	minAPIVersionForEnumValue = map[string]map[string]string{}
)

func Add(name string, kind reflect.Type) {
	t[name] = kind
}

func AddMinAPIVersionForType(name, minAPIVersion string) {
	minAPIVersionForType[name] = minAPIVersion
}

func AddMinAPIVersionForEnumValue(enumName, enumValue, minAPIVersion string) {
	if v, ok := minAPIVersionForEnumValue[enumName]; ok {
		v[enumValue] = minAPIVersion
	} else {
		minAPIVersionForEnumValue[enumName] = map[string]string{
			enumValue: minAPIVersion,
		}
	}
}

type Func func(string) (reflect.Type, bool)

func TypeFunc() Func {
	return func(name string) (reflect.Type, bool) {
		typ, ok := t[name]
		if !ok {
			// The /sdk endpoint does not prefix types with the namespace,
			// but extension endpoints, such as /pbm/sdk do.
			name = strings.TrimPrefix(name, "vim25:")
			typ, ok = t[name]
		}
		return typ, ok
	}
}
