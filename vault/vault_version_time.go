// Copyright IBM Corp. 2026, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import "time"

type VaultVersion struct {
	TimestampInstalled time.Time
	Version            string
	BuildDate          string
}
