// +build !openbsd !arm

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
		testName := "Disk Usage"
		if err != nil {
			Warn(ctx, fmt.Sprintf("Could not obtain partition usage for %s: %v.", partition.Mountpoint, err))
		} else {
			if usage.UsedPercent > 95 {
				SpotWarn(ctx, testName, fmt.Sprintf(partition.Mountpoint+" is %d percent full.", usage.UsedPercent))
				Advise(ctx, "It is recommended to have more than five percent of the partition free.")
			} else if usage.Free < 2<<30 {
				SpotWarn(ctx, testName, partition.Mountpoint+" has %d bytes full.")
				Advise(ctx, "It is recommended to have at least 1 GB of space free per partition.")
			} else {
				SpotOk(ctx, testName, partition.Mountpoint+" usage ok.")
			}
		}

	}
	return nil
}
