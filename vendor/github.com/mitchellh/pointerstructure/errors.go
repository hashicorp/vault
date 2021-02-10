package pointerstructure

import "errors"

var (
	// ErrNotFound is returned if a key in a query can't be found
	ErrNotFound = errors.New("couldn't find key")

	// ErrParse is returned if the query cannot be parsed
	ErrParse = errors.New("first char must be '/'")

	// ErrOutOfRange is returned if a query is referencing a slice
	// or array and the requested index is not in the range [0,len(item))
	ErrOutOfRange = errors.New("out of range")

	// ErrInvalidKind is returned if the item is not a map, slice,
	// array, or struct
	ErrInvalidKind = errors.New("invalid value kind")

	// ErrConvert is returned if an item is not of a requested type
	ErrConvert = errors.New("couldn't convert value")
)
