// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build !linux
// +build !linux

package homedir // import "github.com/ory/dockertest/v3/docker/pkg/homedir"

import (
	"errors"
)

// GetStatic is not needed for non-linux systems.
// (Precisely, it is needed only for glibc-based linux systems.)
func GetStatic() (string, error) {
	return "", errors.New("homedir.GetStatic() is not supported on this system")
}
