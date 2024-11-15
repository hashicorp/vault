//go:build darwin && cgo && !ios
// +build darwin,cgo,!ios

package disk

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include <stdint.h>
#include <CoreFoundation/CoreFoundation.h>
#include "iostat_darwin.h"
*/
import "C"

import (
	"context"

	"github.com/shirou/gopsutil/v3/internal/common"
)

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	var buf [C.NDRIVE]C.DriveStats
	n, err := C.gopsutil_v3_readdrivestat(&buf[0], C.int(len(buf)))
	if err != nil {
		return nil, err
	}
	ret := make(map[string]IOCountersStat, 0)
	for i := 0; i < int(n); i++ {
		d := IOCountersStat{
			ReadBytes:  uint64(buf[i].read),
			WriteBytes: uint64(buf[i].written),
			ReadCount:  uint64(buf[i].nread),
			WriteCount: uint64(buf[i].nwrite),
			ReadTime:   uint64(buf[i].readtime / 1000 / 1000), // note: read/write time are in ns, but we want ms.
			WriteTime:  uint64(buf[i].writetime / 1000 / 1000),
			IoTime:     uint64((buf[i].readtime + buf[i].writetime) / 1000 / 1000),
			Name:       C.GoString(&buf[i].name[0]),
		}
		if len(names) > 0 && !common.StringsHas(names, d.Name) {
			continue
		}

		ret[d.Name] = d
	}
	return ret, nil
}
