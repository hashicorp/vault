// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import "errors"

var (
	// ErrInternal should be used to represent an unexpected error that occurred
	// within the audit system.
	ErrInternal = errors.New("audit system internal error")

	// ErrInvalidParameter should be used to represent an error in which the
	// internal audit system is receiving invalid parameters from other parts of
	// Vault which should have already been validated.
	ErrInvalidParameter = errors.New("invalid internal parameter")

	// ErrExternalOptions should be used to represent an error related to
	// invalid configuration provided to Vault (i.e. by the Vault Operator).
	ErrExternalOptions = errors.New("invalid configuration")
)

// ConvertToExternalError handles converting an audit related error that was generated
// in Vault and should appear as-is in the server logs, to an error that can be
// returned to calling clients (via the API/CLI).
func ConvertToExternalError(err error) error {
	// If the error is an internal error, the contents will have been logged, and
	// we should probably shield the caller from the details.
	if errors.Is(err, ErrInternal) {
		return ErrInternal
	}

	return err
}
