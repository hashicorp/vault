// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

func forwardPkiCertCounts(c *Core, issuedCount uint64, storedCount uint64) bool {
	return false
}
