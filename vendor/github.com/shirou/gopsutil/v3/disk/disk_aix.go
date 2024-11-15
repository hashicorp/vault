//go:build aix
// +build aix

package disk

import (
	"context"
	"errors"
	"strings"

	"github.com/shirou/gopsutil/v3/internal/common"
)

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}

func LabelWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}

// Using lscfg and a device name, we can get the device information
// This is a pure go implementation, and should be moved to disk_aix_nocgo.go
// if a more efficient CGO method is introduced in disk_aix_cgo.go
func SerialNumberWithContext(ctx context.Context, name string) (string, error) {
	// This isn't linux, these aren't actual disk devices
	if strings.HasPrefix(name, "/dev/") {
		return "", errors.New("devices on /dev are not physical disks on aix")
	}
	out, err := invoke.CommandWithContext(ctx, "lscfg", "-vl", name)
	if err != nil {
		return "", err
	}

	ret := ""
	// Kind of inefficient, but it works
	lines := strings.Split(string(out[:]), "\n")
	for line := 1; line < len(lines); line++ {
		v := strings.TrimSpace(lines[line])
		if strings.HasPrefix(v, "Serial Number...............") {
			ret = strings.TrimPrefix(v, "Serial Number...............")
			if ret == "" {
				return "", errors.New("empty serial for disk")
			}
			return ret, nil
		}
	}

	return ret, errors.New("serial entry not found for disk")
}
