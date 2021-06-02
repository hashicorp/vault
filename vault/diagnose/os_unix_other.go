// +build !windows !openbsd,!arm

package diagnose

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

func diskUsage(ctx context.Context) error {
	// Disk usage
	partitions, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	partitionExcludes := []string{"/boot"}
partLoop:
	for _, partition := range partitions {
		for _, exc := range partitionExcludes {
			if strings.HasPrefix(partition.Mountpoint, exc) {
				continue partLoop
			}
		}
		usage, err := disk.Usage(partition.Mountpoint)
		testName := "disk usage"
		if err != nil {
			Warn(ctx, fmt.Sprintf("could not obtain partition usage for %s: %v", partition.Mountpoint, err))
		} else {
			if usage.UsedPercent > 95 {
				SpotWarn(ctx, testName, partition.Mountpoint+" more than 95% full")
			} else if usage.Free < 2<<30 {
				SpotWarn(ctx, testName, partition.Mountpoint+" less than 1GB free")
			} else {
				SpotOk(ctx, testName, partition.Mountpoint+" ok")
			}
		}

	}
	return nil
}
