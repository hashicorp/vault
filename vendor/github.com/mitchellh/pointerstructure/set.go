package pointerstructure

import (
	"fmt"
	"reflect"
)

// Set writes a value v to the pointer p in structure s.
//
// The structures s must have non-zero values set up to this pointer.
// For example, if setting "/bob/0/name", then "/bob/0" must be set already.
//
// The returned value is potentially a new value if this pointer represents
// the root document. Otherwise, the returned value will always be s.
func (p *Pointer) Set(s, v interface{}) (interface{}, error) {
	// if we represent the root doc, return that
	if len(p.Parts) == 0 {
		return v, nil
	}

	// Save the original since this is going to be our return value
	originalS := s

	// Get the parent value
	var err error
	s, err = p.Parent().Get(s)
	if err != nil {
		return nil, err
	}

	// Map for lookup of getter to call for type
	funcMap := map[reflect.Kind]setFunc{
		reflect.Array: p.setSlice,
		reflect.Map:   p.setMap,
		reflect.Slice: p.setSlice,
	}

	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	for val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	f, ok := funcMap[val.Kind()]
	if !ok {
		return nil, fmt.Errorf("set %s: %w: %s", p, ErrInvalidKind, val.Kind())
	}

	result, err := f(originalS, val, reflect.ValueOf(v))
	if err != nil {
		return nil, fmt.Errorf("set %s: %w", p, err)
	}

	return result, nil
}

type setFunc func(interface{}, reflect.Value, reflect.Value) (interface{}, error)

func (p *Pointer) setMap(root interface{}, m, value reflect.Value) (interface{}, error) {
	part := p.Parts[len(p.Parts)-1]
	key, err := coerce(reflect.ValueOf(part), m.Type().Key())
	if err != nil {
		return root, err
	}

	elem, err := coerce(value, m.Type().Elem())
	if err != nil {
		return root, err
	}

	// Set the key
	m.SetMapIndex(key, elem)
	return root, nil
}

func (p *Pointer) setSlice(root interface{}, s, value reflect.Value) (interface{}, error) {
	// Coerce the value, we'll need that no matter what
	value, err := coerce(value, s.Type().Elem())
	if err != nil {
		return root, err
	}

	// If the part is the special "-", that means to append it (RFC6901 4.)
	part := p.Parts[len(p.Parts)-1]
	if part == "-" {
		return p.setSliceAppend(root, s, value)
	}

	// Coerce the key to an int
	idxVal, err := coerce(reflect.ValueOf(part), reflect.TypeOf(42))
	if err != nil {
		return root, err
	}
	idx := int(idxVal.Int())

	// Verify we're within bounds
	if idx < 0 || idx >= s.Len() {
		return root, fmt.Errorf(
			"index %d is %w (length = %d)", idx, ErrOutOfRange, s.Len())
	}

	// Set the key
	s.Index(idx).Set(value)
	return root, nil
}

func (p *Pointer) setSliceAppend(root interface{}, s, value reflect.Value) (interface{}, error) {
	// Coerce the value, we'll need that no matter what. This should
	// be a no-op since we expect it to be done already, but there is
	// a fast-path check for that in coerce so do it anyways.
	value, err := coerce(value, s.Type().Elem())
	if err != nil {
		return root, err
	}

	// We can assume "s" is the parent of pointer value. We need to actually
	// write s back because Append can return a new slice.
	return p.Parent().Set(root, reflect.Append(s, value).Interface())
}
