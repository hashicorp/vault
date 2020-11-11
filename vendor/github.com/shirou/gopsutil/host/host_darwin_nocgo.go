// +build darwin
// +build !cgo

package host

import (
	"context"

	"github.com/shirou/gopsutil/internal/common"
)

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}
