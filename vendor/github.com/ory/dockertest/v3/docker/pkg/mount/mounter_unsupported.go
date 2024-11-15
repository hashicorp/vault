// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build (!linux && !freebsd) || (freebsd && !cgo)
// +build !linux,!freebsd freebsd,!cgo

package mount // import "github.com/ory/dockertest/v3/docker/pkg/mount"

func mount(device, target, mType string, flag uintptr, data string) error {
	panic("Not implemented")
}

func unmount(target string, flag int) error {
	panic("Not implemented")
}
