// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

// MFAMethodInfo contains the information about an MFA method
type MFAMethodInfo struct {
	MethodID    string
	MethodType  string
	UsePasscode bool
}
