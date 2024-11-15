// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fileutils // import "github.com/ory/dockertest/v3/docker/pkg/fileutils"

// GetTotalUsedFds Returns the number of used File Descriptors. Not supported
// on Windows.
func GetTotalUsedFds() int {
	return -1
}
