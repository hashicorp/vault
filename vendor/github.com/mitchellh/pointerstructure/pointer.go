// Package pointerstructure provides functions for identifying a specific
// value within any Go structure using a string syntax.
//
// The syntax used is based on JSON Pointer (RFC 6901).
package pointerstructure

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Pointer represents a pointer to a specific value. You can construct
// a pointer manually or use Parse.
type Pointer struct {
	// Parts are the pointer parts. No escape codes are processed here.
	// The values are expected to be exact. If you have escape codes, use
	// the Parse functions.
	Parts []string
}

// Get reads the value at the given pointer.
//
// This is a shorthand for calling Parse on the pointer and then calling Get
// on that result. An error will be returned if the value cannot be found or
// there is an error with the format of pointer.
func Get(value interface{}, pointer string) (interface{}, error) {
	p, err := Parse(pointer)
	if err != nil {
		return nil, err
	}

	return p.Get(value)
}

// Set sets the value at the given pointer.
//
// This is a shorthand for calling Parse on the pointer and then calling Set
// on that result. An error will be returned if the value cannot be found or
// there is an error with the format of pointer.
//
// Set returns the complete document, which might change if the pointer value
// points to the root ("").
func Set(doc interface{}, pointer string, value interface{}) (interface{}, error) {
	p, err := Parse(pointer)
	if err != nil {
		return nil, err
	}

	return p.Set(doc, value)
}

// String returns the string value that can be sent back to Parse to get
// the same Pointer result.
func (p *Pointer) String() string {
	if len(p.Parts) == 0 {
		return ""
	}

	// Copy the parts so we can convert back the escapes
	result := make([]string, len(p.Parts))
	copy(result, p.Parts)
	for i, p := range p.Parts {
		result[i] = strings.Replace(
			strings.Replace(p, "~", "~0", -1), "/", "~1", -1)

	}

	return "/" + strings.Join(result, "/")
}

// Parent returns a pointer to the parent element of this pointer.
//
// If Pointer represents the root (empty parts), a pointer representing
// the root is returned. Therefore, to check for the root, IsRoot() should be
// called.
func (p *Pointer) Parent() *Pointer {
	// If this is root, then we just return a new root pointer. We allocate
	// a new one though so this can still be modified.
	if p.IsRoot() {
		return &Pointer{}
	}

	parts := make([]string, len(p.Parts)-1)
	copy(parts, p.Parts[:len(p.Parts)-1])
	return &Pointer{
		Parts: parts,
	}
}

// IsRoot returns true if this pointer represents the root document.
func (p *Pointer) IsRoot() bool {
	return len(p.Parts) == 0
}

// coerce is a helper to coerce a value to a specific type if it must
// and if its possible. If it isn't possible, an error is returned.
func coerce(value reflect.Value, to reflect.Type) (reflect.Value, error) {
	// If the value is already assignable to the type, then let it go
	if value.Type().AssignableTo(to) {
		return value, nil
	}

	// If a direct conversion is possible, do that
	if value.Type().ConvertibleTo(to) {
		return value.Convert(to), nil
	}

	// Create a new value to hold our result
	result := reflect.New(to)

	// Decode
	if err := mapstructure.WeakDecode(value.Interface(), result.Interface()); err != nil {
		return result, fmt.Errorf(
			"couldn't convert value %#v to type %s",
			value.Interface(), to.String())
	}

	// We need to indirect the value since reflect.New always creates a pointer
	return reflect.Indirect(result), nil
}
