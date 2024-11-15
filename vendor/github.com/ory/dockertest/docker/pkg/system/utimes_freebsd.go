package system // import "github.com/ory/dockertest/docker/pkg/system"

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// LUtimesNano is used to change access and modification time of the specified path.
// It's used for symbol link file because unix.UtimesNano doesn't support a NOFOLLOW flag atm.
func LUtimesNano(path string, ts []syscall.Timespec) error {
	var _path *byte
	_path, err := unix.BytePtrFromString(path)
	if err != nil {
		return err
	}

	if _, _, err := unix.Syscall(unix.SYS_LUTIMES, uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&ts[0])), 0); err != 0 && err != unix.ENOSYS {
		return err
	}

	return nil
}
