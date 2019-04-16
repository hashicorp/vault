package pointerstructure

import (
	"fmt"
	"reflect"
)

// Get reads the value out of the total value v.
func (p *Pointer) Get(v interface{}) (interface{}, error) {
	// fast-path the empty address case to avoid reflect.ValueOf below
	if len(p.Parts) == 0 {
		return v, nil
	}

	// Map for lookup of getter to call for type
	funcMap := map[reflect.Kind]func(string, reflect.Value) (reflect.Value, error){
		reflect.Array:  p.getSlice,
		reflect.Map:    p.getMap,
		reflect.Slice:  p.getSlice,
		reflect.Struct: p.getStruct,
	}

	currentVal := reflect.ValueOf(v)
	for i, part := range p.Parts {
		for currentVal.Kind() == reflect.Interface {
			currentVal = currentVal.Elem()
		}

		for currentVal.Kind() == reflect.Ptr {
			currentVal = reflect.Indirect(currentVal)
		}

		f, ok := funcMap[currentVal.Kind()]
		if !ok {
			return nil, fmt.Errorf(
				"%s: at part %d, invalid value kind: %s", p, i, currentVal.Kind())
		}

		var err error
		currentVal, err = f(part, currentVal)
		if err != nil {
			return nil, fmt.Errorf("%s at part %d: %s", p, i, err)
		}
	}

	return currentVal.Interface(), nil
}

func (p *Pointer) getMap(part string, m reflect.Value) (reflect.Value, error) {
	var zeroValue reflect.Value

	// Coerce the string part to the correct key type
	key, err := coerce(reflect.ValueOf(part), m.Type().Key())
	if err != nil {
		return zeroValue, err
	}

	// Verify that the key exists
	found := false
	for _, k := range m.MapKeys() {
		if k.Interface() == key.Interface() {
			found = true
			break
		}
	}
	if !found {
		return zeroValue, fmt.Errorf("couldn't find key %#v", key.Interface())
	}

	// Get the key
	return m.MapIndex(key), nil
}

func (p *Pointer) getSlice(part string, v reflect.Value) (reflect.Value, error) {
	var zeroValue reflect.Value

	// Coerce the key to an int
	idxVal, err := coerce(reflect.ValueOf(part), reflect.TypeOf(42))
	if err != nil {
		return zeroValue, err
	}
	idx := int(idxVal.Int())

	// Verify we're within bounds
	if idx < 0 || idx >= v.Len() {
		return zeroValue, fmt.Errorf(
			"index %d is out of range (length = %d)", idx, v.Len())
	}

	// Get the key
	return v.Index(idx), nil
}

func (p *Pointer) getStruct(part string, m reflect.Value) (reflect.Value, error) {
	return m.FieldByName(part), nil
}
