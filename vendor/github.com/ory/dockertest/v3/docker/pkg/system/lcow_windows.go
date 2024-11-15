// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

// LCOWSupported returns true if Linux containers on Windows are supported.
func LCOWSupported() bool {
	return lcowSupported
}
