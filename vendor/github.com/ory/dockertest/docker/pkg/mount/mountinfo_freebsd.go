package mount // import "github.com/ory/dockertest/docker/pkg/mount"

/*
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Parse /proc/self/mountinfo because comparing Dev and ino does not work from
// bind mounts.
func parseMountTable() ([]*Info, error) {
	var rawEntries *C.struct_statfs

	count := int(C.getmntinfo(&rawEntries, C.MNT_WAIT))
	if count == 0 {
		return nil, fmt.Errorf("Failed to call getmntinfo")
	}

	var entries []C.struct_statfs
	header := (*reflect.SliceHeader)(unsafe.Pointer(&entries))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(rawEntries))

	var out []*Info
	for _, entry := range entries {
		var mountinfo Info
		mountinfo.Mountpoint = C.GoString(&entry.f_mntonname[0])
		mountinfo.Source = C.GoString(&entry.f_mntfromname[0])
		mountinfo.Fstype = C.GoString(&entry.f_fstypename[0])
		out = append(out, &mountinfo)
	}
	return out, nil
}
