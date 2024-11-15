//go:build netbsd
// +build netbsd

package host

import (
	"context"
	"strings"

	"github.com/shirou/gopsutil/v3/internal/common"
	"golang.org/x/sys/unix"
)

func HostIDWithContext(ctx context.Context) (string, error) {
	return "", common.ErrNotImplementedError
}

func numProcs(ctx context.Context) (uint64, error) {
	return 0, common.ErrNotImplementedError
}

func PlatformInformationWithContext(ctx context.Context) (string, string, string, error) {
	platform := ""
	family := ""
	version := ""

	p, err := unix.Sysctl("kern.ostype")
	if err == nil {
		platform = strings.ToLower(p)
	}
	v, err := unix.Sysctl("kern.osrelease")
	if err == nil {
		version = strings.ToLower(v)
	}

	return platform, family, version, nil
}

func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	return "", "", common.ErrNotImplementedError
}

func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	var ret []UserStat
	return ret, common.ErrNotImplementedError
}

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}

func KernelVersionWithContext(ctx context.Context) (string, error) {
	_, _, version, err := PlatformInformationWithContext(ctx)
	return version, err
}
