// +build go1.12

package fmtsort

import "reflect"

const brokenNaNs = false

func mapElems(mapValue reflect.Value) ([]reflect.Value, []reflect.Value) {
	key := make([]reflect.Value, mapValue.Len())
	value := make([]reflect.Value, len(key))
	iter := mapValue.MapRange()
	for i := 0; iter.Next(); i++ {
		key[i] = iter.Key()
		value[i] = iter.Value()
	}
	return key, value
}
