package api

import (
	"context"
	"errors"
	"time"

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

	// Parse timestamp separately since WeakDecode can't handle this field.
	timestampRaw := secret.Data["timestamp"].(string)
	timestamp, err := time.Parse(time.RFC3339, timestampRaw)
	if err != nil {
		return nil, err
	}
	delete(secret.Data, "timestamp")

	var result HostInfoResponse
	err = mapstructure.WeakDecode(secret.Data, &result)
	if err != nil {
		return nil, err
	}
	result.Timestamp = timestamp

	return &result, err
}

type HostInfoResponse struct {
	Timestamp time.Time              `json:"timestamp"`
	CPU       []cpu.InfoStat         `json:"cpu"`
	Disk      []*disk.UsageStat      `json:"disk"`
	Host      *host.InfoStat         `json:"host"`
	Memory    *mem.VirtualMemoryStat `json:"memory"`
}
