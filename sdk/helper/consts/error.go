// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package consts

import (
	"errors"
	"strings"
)

var (
	// ErrSealed is returned if an operation is performed on a sealed barrier.
	// No operation is expected to succeed before unsealing
	ErrSealed = errors.New("Vault is sealed")

	// ErrAPILocked is returned if an operation is performed when the API is
	// locked for the request namespace.
	ErrAPILocked = errors.New("API access to this namespace has been locked by an administrator")

	// ErrStandby is returned if an operation is performed on a standby Vault.
	// No operation is expected to succeed until active.
	ErrStandby = errors.New("Vault is in standby mode")

	// ErrPathContainsParentReferences is returned when a path contains parent
	// references.
	ErrPathContainsParentReferences = errors.New("path cannot contain parent references")

	// ErrInvalidWrappingToken is returned when checking for the validity of
	// a wrapping token that turns out to be invalid.
	ErrInvalidWrappingToken = errors.New("wrapping token is not valid or does not exist")

	// ErrOverloaded indicates the Vault server is at capacity.
	ErrOverloaded = errors.New("overloaded, try again later")
)

// PathContainsParentReferences checks whether a path contains ".." as an
// actual parent directory reference (i.e., a complete path segment), as
// opposed to ".." appearing as a substring of a longer segment like "...".
//
// For example:
//   - "foo/../bar"  => true  (parent reference)
//   - "foo/.."      => true  (parent reference)
//   - "../foo"      => true  (parent reference)
//   - ".."          => true  (parent reference)
//   - "foo/.../bar" => false (three dots, not a parent reference)
//   - "test_..."    => false (dots are part of a longer name)
func PathContainsParentReferences(path string) bool {
	for _, segment := range strings.Split(path, "/") {
		if segment == ".." {
			return true
		}
	}
	return false
}
