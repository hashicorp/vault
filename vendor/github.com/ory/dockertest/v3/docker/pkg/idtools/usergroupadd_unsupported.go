// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !linux
// +build !linux

package idtools // import "github.com/ory/dockertest/v3/docker/pkg/idtools"

import "fmt"

// AddNamespaceRangesUser takes a name and finds an unused uid, gid pair
// and calls the appropriate helper function to add the group and then
// the user to the group in /etc/group and /etc/passwd respectively.
func AddNamespaceRangesUser(name string) (int, int, error) {
	return -1, -1, fmt.Errorf("No support for adding users or groups on this OS")
}
