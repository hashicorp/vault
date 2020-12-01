// +build !windows,!linux,!freebsd freebsd,!cgo

package mount // import "github.com/ory/dockertest/v3/docker/pkg/mount"

import (
	"fmt"
	"runtime"
)

func parseMountTable() ([]*Info, error) {
	return nil, fmt.Errorf("mount.parseMountTable is not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
}
