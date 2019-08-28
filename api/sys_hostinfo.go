package api

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
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

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result HostInfoResponse
	err = mapstructure.WeakDecode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

type HostInfoResponse struct {
	Timestamp string                 `json:"timestamp" mapstructure:"-"`
	CPU       []cpu.InfoStat         `json:"cpu"`
	Disk      []*disk.UsageStat      `json:"disk"`
	Host      *host.InfoStat         `json:"host"`
	Memory    *mem.VirtualMemoryStat `json:"memory"`
}
