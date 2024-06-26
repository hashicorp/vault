// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokenhelper

// TokenHelper is an interface that contains basic operations that must be
// implemented by a token helper
type TokenHelper interface {
	// Path displays a method-specific path; for the internal helper this
	// is the location of the token stored on disk; for the external helper
	// this is the location of the binary being invoked
	Path() string
	Erase() error
	Get() (string, error)
	Store(string) error
}
