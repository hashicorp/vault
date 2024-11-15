// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

// lcowSupported determines if Linux Containers on Windows are supported.
var lcowSupported = false

// InitLCOW sets whether LCOW is supported or not
func InitLCOW(experimental bool) {
	v := GetOSVersion()
	if experimental && v.Build >= 16299 {
		lcowSupported = true
	}
}
