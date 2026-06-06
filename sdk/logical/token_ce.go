// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package logical

type (
	EntToken struct{}
)

// IsStorageBacked reports whether this token has backing token storage.
// Community edition tokens always use token storage.
func (te *TokenEntry) IsStorageBacked() bool {
	return te != nil
}
