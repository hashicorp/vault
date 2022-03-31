//go:build windows

package osutil

import (
	"io/fs"
)

func FileUidMatch(info fs.FileInfo, path string, uid int) error {
	return nil
}
