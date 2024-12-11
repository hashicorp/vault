// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"time"

	"github.com/hashicorp/vault/version"
)

func init() {
	// The BuildDate is set as part of the build process in CI so we need to
	// initialize it for testing. By setting it to now minus one year we
	// provide some headroom to ensure that test license expiration (for enterprise)
	// does not exceed the BuildDate as that is invalid.
	if version.BuildDate == "" {
		version.BuildDate = time.Now().UTC().AddDate(-1, 0, 0).Format(time.RFC3339)
	}
}
