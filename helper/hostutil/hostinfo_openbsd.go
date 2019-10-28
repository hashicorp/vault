// +build openbsd

package hostutil

import (
	"fmt"
	"time"
)

type HostInfo struct {
	Timestamp time.Time     `json:"timestamp"`
	CPU       []interface{} `json:"cpu"`
	CPUTimes  []interface{} `json:"cpu_times"`
	Disk      []interface{} `json:"disk"`
	Host      interface{}   `json:"host"`
	Memory    interface{}   `json:"memory"`
}

func CollectHostInfo() (*HostInfo, error) {
	return nil, fmt.Errorf("host info not supported on this platform")
}
