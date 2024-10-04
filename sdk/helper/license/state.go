// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package license

import "time"

type LicenseState struct {
	State      string
	ExpiryTime time.Time
	Terminated bool
}
