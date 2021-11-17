package crypto

import "errors"

var (
	// ErrInvalidKey means the key argument passed to SigningMethod.Verify
	// was not the correct type.
	ErrInvalidKey = errors.New("key is invalid")
)
