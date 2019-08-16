package api

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func (c *Sys) HostInfo() (*HostInfoResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/host-info")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result HostInfoResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type HostInfoResponse struct {
	CollectionTime time.Time              `json:"collection_time"`
	CPU            []cpu.InfoStat         `json:"cpu"`
	Disk           *disk.UsageStat        `json:"disk"`
	Host           *host.InfoStat         `json:"host"`
	Memory         *mem.VirtualMemoryStat `json:"memory"`
}
