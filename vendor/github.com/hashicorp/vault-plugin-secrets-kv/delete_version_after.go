// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kv

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
)

// deletionTime returns the time of creation plus the duration of the
// minimum non-zero value of mount or meta. If mount and meta are zero,
// false is returned.
func deletionTime(creation time.Time, mount, meta time.Duration) (time.Time, bool) {
	if mount == 0 && meta == 0 {
		return time.Time{}, false
	}
	var min time.Duration
	if meta != 0 {
		min = meta
	}
	if (mount != 0 && mount < min) || min == 0 {
		min = mount
	}
	return creation.Add(min), true
}

type deleteVersionAfterGetter interface {
	GetDeleteVersionAfter() *duration.Duration
}

func deleteVersionAfter(v deleteVersionAfterGetter) time.Duration {
	if v.GetDeleteVersionAfter() == nil {
		return time.Duration(0)
	}
	dva, err := ptypes.Duration(v.GetDeleteVersionAfter())
	if err != nil {
		return time.Duration(0)
	}
	return dva
}

const (
	disabled time.Duration = -1 * time.Second
)

// IsDeleteVersionAfterDisabled returns true if DeleteVersionAfter is
// disabled.
func (c *Configuration) IsDeleteVersionAfterDisabled() bool {
	if deleteVersionAfter(c) == disabled {
		return true
	}
	return false
}

// DisableDeleteVersionAfter disables DeleteVersionAfter.
func (c *Configuration) DisableDeleteVersionAfter() {
	c.DeleteVersionAfter = ptypes.DurationProto(disabled)
}

// ResetDeleteVersionAfter resets the DeleteVersionAfter to the default
// value.
func (c *Configuration) ResetDeleteVersionAfter() {
	c.DeleteVersionAfter = nil
}
