// +build darwin
// +build !cgo

package disk

import (
	"context"

	"github.com/shirou/gopsutil/internal/common"
)

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}
