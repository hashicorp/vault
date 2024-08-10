// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package license

import "time"

type LicenseState struct {
	State      string
	ExpiryTime time.Time
	Terminated bool
}
