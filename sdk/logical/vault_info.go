// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"time"
)

type VaultInfo struct {
	BuildDate         time.Time
	BuiltinPublicKeys []string
}
