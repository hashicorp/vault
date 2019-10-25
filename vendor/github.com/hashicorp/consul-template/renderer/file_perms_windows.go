//+build windows

package renderer

import "os"

func preserveFilePermissions(path string, fileInfo os.FileInfo) error {
	return nil
}
