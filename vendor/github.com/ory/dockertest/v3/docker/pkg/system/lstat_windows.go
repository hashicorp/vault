// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package system // import "github.com/ory/dockertest/v3/docker/pkg/system"

import "os"

// Lstat calls os.Lstat to get a fileinfo interface back.
// This is then copied into our own locally defined structure.
func Lstat(path string) (*StatT, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	return fromStatT(&fi)
}
