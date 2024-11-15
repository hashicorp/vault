// +build !windows

package atomic

import (
	"os"
)

// ReplaceFile atomically replaces the destination file or directory with the
// source.  It is guaranteed to either replace the target file entirely, or not
// change either file.
func ReplaceFile(source, destination string) error {
	return os.Rename(source, destination)
}
