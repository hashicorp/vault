// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import "time"

type VaultVersion struct {
	TimestampInstalled time.Time
	Version            string
	BuildDate          string
}
