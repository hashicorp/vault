// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package errutil

// UserError represents an error generated due to invalid user input
type UserError struct {
	Err string
}

func (e UserError) Error() string {
	return e.Err
}

// InternalError represents an error generated internally,
// presumably not due to invalid user input
type InternalError struct {
	Err string
}

func (e InternalError) Error() string {
	return e.Err
}
