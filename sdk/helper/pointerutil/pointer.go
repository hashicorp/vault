// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pointerutil

import (
	"os"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

// StringPtr returns a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a boolean value
func BoolPtr(b bool) *bool {
	return &b
}

// TimeDurationPtr returns a pointer to a time duration value
func TimeDurationPtr(duration string) *time.Duration {
	d, _ := parseutil.ParseDurationSecond(duration)

	return &d
}

// FileModePtr returns a pointer to the given os.FileMode
func FileModePtr(o os.FileMode) *os.FileMode {
	return &o
}

// Int64Ptr returns a pointer to an int64 value
func Int64Ptr(i int64) *int64 {
	return &i
}
