// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logging

import (
	"os"
	"time"
)

func (l *LogFile) createTime(stat os.FileInfo) time.Time {
	// Use `ModTime` as an approximation if the exact create time is not
	// available.
	// On Windows, the file create time is not updated after the active log
	// rotates, so use `ModTime` as an approximation as well.
	return stat.ModTime()
}
