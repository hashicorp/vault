// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import "path/filepath"

const (
	// umaskBaseOctal and umaskBitSize are used when parsing umask values.
	umaskBaseOctal = 8
	umaskBitSize   = 32
)

// seenBlockName tracks the array index and source line of a previously seen
// block name so duplicate-name errors can reference both locations.
type seenBlockName struct {
	idx  int
	line int
}

// normalizeDestinationPath converts path to an absolute, cleaned path and
// resolves symlinks when possible. If Abs or EvalSymlinks fails, it falls back
// to returning the best cleaned path available.
func normalizeDestinationPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	absPath = filepath.Clean(absPath)
	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		return resolved
	}
	return absPath
}
