package hostutil

import (
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// VirutalMemoryStat holds commonly used memory measurements. We must have a
// local type here in order to avoid building the gopsutil library on certain
// arch types.
type VirtualMemoryStat struct {
	*mem.VirtualMemoryStat

	// A subset of JSON struct tags in mem.VirtualMemoryStat were changed in a backwards
	// incompatible way in the v3 release of gopsutil. The fields below are copied from
	// v2 of gopsutil to maintain backwards compatibility in the Vault host-info API.
	//
	// The following details the changed JSON struct tags between v2 and v3:
	// https://github.com/shirou/gopsutil/blob/master/_tools/v3migration/v3migration.sh#L61
	CommitLimit    uint64 `json:"commitlimit"`
	CommittedAS    uint64 `json:"committedas"`
	HighFree       uint64 `json:"highfree"`
	HighTotal      uint64 `json:"hightotal"`
	HugePagesFree  uint64 `json:"hugepagesfree"`
	HugePageSize   uint64 `json:"hugepagesize"`
	HugePagesTotal uint64 `json:"hugepagestotal"`
	LowFree        uint64 `json:"lowfree"`
	LowTotal       uint64 `json:"lowtotal"`
	PageTables     uint64 `json:"pagetables"`
	SwapCached     uint64 `json:"swapcached"`
	SwapFree       uint64 `json:"swapfree"`
	SwapTotal      uint64 `json:"swaptotal"`
	VMallocChunk   uint64 `json:"vmallocchunk"`
	VMallocTotal   uint64 `json:"vmalloctotal"`
	VMallocUsed    uint64 `json:"vmallocused"`
	Writeback      uint64 `json:"writeback"`
	WritebackTmp   uint64 `json:"writebacktmp"`
}

// A HostInfoStat describes the host status.
type HostInfoStat struct {
	*host.InfoStat

	// A subset of JSON struct tags in host.InfoStat were changed in a backwards
	// incompatible way in the v3 release of gopsutil. The fields below are copied from
	// v2 of gopsutil to maintain backwards compatibility in the Vault host-info API.
	//
	// The following details the changed JSON struct tags between v2 and v3:
	// https://github.com/shirou/gopsutil/blob/master/_tools/v3migration/v3migration.sh#L72
	HostID string `json:"hostid"`
}
