// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

// Umask is not supported on the windows platform.
func Umask(newmask int) (oldmask int, err error) {
	// should not be called on cli code path
	return 0, ErrNotSupportedPlatform
}
