// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldap

import "errors"

var (
	// ErrUnknown is an unknown/undefined error
	ErrUnknown = errors.New("unknown")

	// ErrInvalidParameter is an invalid parameter error
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrInternal is an internal error
	ErrInternal = errors.New("internal error")
)
