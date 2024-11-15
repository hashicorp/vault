package pointerstructure

import (
	"fmt"
	"reflect"
)

// Delete deletes the value specified by the pointer p in structure s.
//
// When deleting a slice index, all other elements will be shifted to
// the left. This is specified in RFC6902 (JSON Patch) and not RFC6901 since
// RFC6901 doesn't specify operations on pointers. If you don't want to
// shift elements, you should use Set to set the slice index to the zero value.
//
// The structures s must have non-zero values set up to this pointer.
// For example, if deleting "/bob/0/name", then "/bob/0" must be set already.
//
// The returned value is potentially a new value if this pointer represents
// the root document. Otherwise, the returned value will always be s.
func (p *Pointer) Delete(s interface{}) (interface{}, error) {
	// if we represent the root doc, we've deleted everything
	if len(p.Parts) == 0 {
		return nil, nil
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
	funcMap := map[reflect.Kind]deleteFunc{
		reflect.Array: p.deleteSlice,
		reflect.Map:   p.deleteMap,
		reflect.Slice: p.deleteSlice,
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
		return nil, fmt.Errorf("delete %s: %w: %s", p, ErrInvalidKind, val.Kind())
	}

	result, err := f(originalS, val)
	if err != nil {
		return nil, fmt.Errorf("delete %s: %s", p, err)
	}

	return result, nil
}

type deleteFunc func(interface{}, reflect.Value) (interface{}, error)

func (p *Pointer) deleteMap(root interface{}, m reflect.Value) (interface{}, error) {
	part := p.Parts[len(p.Parts)-1]
	key, err := coerce(reflect.ValueOf(part), m.Type().Key())
	if err != nil {
		return root, err
	}

	// Delete the key
	var elem reflect.Value
	m.SetMapIndex(key, elem)
	return root, nil
}

func (p *Pointer) deleteSlice(root interface{}, s reflect.Value) (interface{}, error) {
	// Coerce the key to an int
	part := p.Parts[len(p.Parts)-1]
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

	// Mimicing the following with reflection to do this:
	//
	// copy(a[i:], a[i+1:])
	// a[len(a)-1] = nil // or the zero value of T
	// a = a[:len(a)-1]

	// copy(a[i:], a[i+1:])
	reflect.Copy(s.Slice(idx, s.Len()), s.Slice(idx+1, s.Len()))

	// a[len(a)-1] = nil // or the zero value of T
	s.Index(s.Len() - 1).Set(reflect.Zero(s.Type().Elem()))

	// a = a[:len(a)-1]
	s = s.Slice(0, s.Len()-1)

	// set the slice back on the parent
	return p.Parent().Set(root, s.Interface())
}
