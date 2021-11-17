package pointerstructure

import (
	"fmt"
	"reflect"
	"strings"
)

// Get reads the value out of the total value v.
//
// For struct values a `pointer:"<name>"` tag on the struct's
// fields may be used to override that field's name for lookup purposes.
// Alternatively the tag name used can be overridden in the `Config`.
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
				"%s: at part %d, %w: %s", p, i, ErrInvalidKind, currentVal.Kind())
		}

		var err error
		currentVal, err = f(part, currentVal)
		if err != nil {
			return nil, fmt.Errorf("%s at part %d: %w", p, i, err)
		}
		if p.Config.ValueTransformationHook != nil {
			currentVal = p.Config.ValueTransformationHook(currentVal)
			if currentVal == reflect.ValueOf(nil) {
				return nil, fmt.Errorf("%s at part %d: ValueTransformationHook returned the value of a nil interface", p, i)
			}
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
		return zeroValue, fmt.Errorf("%w %#v", ErrNotFound, key.Interface())
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
			"index %d is %w (length = %d)", idx, ErrOutOfRange, v.Len())
	}

	// Get the key
	return v.Index(idx), nil
}

func (p *Pointer) getStruct(part string, m reflect.Value) (reflect.Value, error) {
	var foundField reflect.Value
	var found bool
	var ignored bool
	typ := m.Type()

	tagName := p.Config.TagName
	if tagName == "" {
		tagName = "pointer"
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.PkgPath != "" {
			// this is an unexported field so ignore it
			continue
		}

		fieldTag := field.Tag.Get(tagName)

		if fieldTag != "" {
			if idx := strings.Index(fieldTag, ","); idx != -1 {
				fieldTag = fieldTag[0:idx]
			}

			if strings.Contains(fieldTag, "|") {
				// should this panic instead?
				return foundField, fmt.Errorf("pointer struct tag cannot contain the '|' character")
			}

			if fieldTag == "-" {
				// we should ignore this field but cannot immediately return because its possible another
				// field has a tag that would allow it to assume this ones name.

				if field.Name == part {
					found = true
					ignored = true
				}
				continue
			} else if fieldTag == part {
				// we can go ahead and return now as the tag is enough to
				// indicate that this is the correct field
				return m.Field(i), nil
			}
		} else if field.Name == part {
			foundField = m.Field(i)
			found = true
		}
	}

	if !found {
		return reflect.Value{}, fmt.Errorf("couldn't find struct field with name %q", part)
	}

	if ignored {
		return reflect.Value{}, fmt.Errorf("struct field %q is ignored and cannot be used", part)
	}

	return foundField, nil
}
