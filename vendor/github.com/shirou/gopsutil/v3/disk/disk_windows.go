//go:build windows
// +build windows

package disk

import (
	"bytes"
	"context"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/shirou/gopsutil/v3/internal/common"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	procGetDiskFreeSpaceExW     = common.Modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetLogicalDriveStringsW = common.Modkernel32.NewProc("GetLogicalDriveStringsW")
	procGetDriveType            = common.Modkernel32.NewProc("GetDriveTypeW")
	procGetVolumeInformation    = common.Modkernel32.NewProc("GetVolumeInformationW")
)

var (
	fileFileCompression = int64(16)     // 0x00000010
	fileReadOnlyVolume  = int64(524288) // 0x00080000
)

// diskPerformance is an equivalent representation of DISK_PERFORMANCE in the Windows API.
// https://docs.microsoft.com/fr-fr/windows/win32/api/winioctl/ns-winioctl-disk_performance
type diskPerformance struct {
	BytesRead           int64
	BytesWritten        int64
	ReadTime            int64
	WriteTime           int64
	IdleTime            int64
	ReadCount           uint32
	WriteCount          uint32
	QueueDepth          uint32
	SplitCount          uint32
	QueryTime           int64
	StorageDeviceNumber uint32
	StorageManagerName  [8]uint16
	alignmentPadding    uint32 // necessary for 32bit support, see https://github.com/elastic/beats/pull/16553
}

func init() {
	// enable disk performance counters on Windows Server editions (needs to run as admin)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\PartMgr`, registry.SET_VALUE)
	if err == nil {
		key.SetDWordValue("EnableCounterForIoctl", 1)
		key.Close()
	}
}

func UsageWithContext(ctx context.Context, path string) (*UsageStat, error) {
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	diskret, _, err := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if diskret == 0 {
		return nil, err
	}
	ret := &UsageStat{
		Path:        path,
		Total:       uint64(lpTotalNumberOfBytes),
		Free:        uint64(lpTotalNumberOfFreeBytes),
		Used:        uint64(lpTotalNumberOfBytes) - uint64(lpTotalNumberOfFreeBytes),
		UsedPercent: (float64(lpTotalNumberOfBytes) - float64(lpTotalNumberOfFreeBytes)) / float64(lpTotalNumberOfBytes) * 100,
		// InodesTotal: 0,
		// InodesFree: 0,
		// InodesUsed: 0,
		// InodesUsedPercent: 0,
	}
	return ret, nil
}

// PartitionsWithContext returns disk partitions.
// Since GetVolumeInformation doesn't have a timeout, this method uses context to set deadline by users.
func PartitionsWithContext(ctx context.Context, all bool) ([]PartitionStat, error) {
	warnings := Warnings{
		Verbose: true,
	}

	var errLogicalDrives error
	retChan := make(chan PartitionStat)
	quitChan := make(chan struct{})
	defer close(quitChan)

	getPartitions := func() {
		defer close(retChan)

		lpBuffer := make([]byte, 254)

		diskret, _, err := procGetLogicalDriveStringsW.Call(
			uintptr(len(lpBuffer)),
			uintptr(unsafe.Pointer(&lpBuffer[0])))
		if diskret == 0 {
			errLogicalDrives = err
			return
		}
		for _, v := range lpBuffer {
			if v >= 65 && v <= 90 {
				path := string(v) + ":"
				typepath, _ := windows.UTF16PtrFromString(path)
				typeret, _, _ := procGetDriveType.Call(uintptr(unsafe.Pointer(typepath)))
				if typeret == 0 {
					err := windows.GetLastError()
					warnings.Add(err)
					continue
				}
				// 2: DRIVE_REMOVABLE 3: DRIVE_FIXED 4: DRIVE_REMOTE 5: DRIVE_CDROM

				if typeret == 2 || typeret == 3 || typeret == 4 || typeret == 5 {
					lpVolumeNameBuffer := make([]byte, 256)
					lpVolumeSerialNumber := int64(0)
					lpMaximumComponentLength := int64(0)
					lpFileSystemFlags := int64(0)
					lpFileSystemNameBuffer := make([]byte, 256)
					volpath, _ := windows.UTF16PtrFromString(string(v) + ":/")
					driveret, _, err := procGetVolumeInformation.Call(
						uintptr(unsafe.Pointer(volpath)),
						uintptr(unsafe.Pointer(&lpVolumeNameBuffer[0])),
						uintptr(len(lpVolumeNameBuffer)),
						uintptr(unsafe.Pointer(&lpVolumeSerialNumber)),
						uintptr(unsafe.Pointer(&lpMaximumComponentLength)),
						uintptr(unsafe.Pointer(&lpFileSystemFlags)),
						uintptr(unsafe.Pointer(&lpFileSystemNameBuffer[0])),
						uintptr(len(lpFileSystemNameBuffer)))
					if driveret == 0 {
						if typeret == 5 || typeret == 2 {
							continue // device is not ready will happen if there is no disk in the drive
						}
						warnings.Add(err)
						continue
					}
					opts := []string{"rw"}
					if lpFileSystemFlags&fileReadOnlyVolume != 0 {
						opts = []string{"ro"}
					}
					if lpFileSystemFlags&fileFileCompression != 0 {
						opts = append(opts, "compress")
					}

					select {
					case retChan <- PartitionStat{
						Mountpoint: path,
						Device:     path,
						Fstype:     string(bytes.ReplaceAll(lpFileSystemNameBuffer, []byte("\x00"), []byte(""))),
						Opts:       opts,
					}:
					case <-quitChan:
						return
					}
				}
			}
		}
	}

	go getPartitions()

	var ret []PartitionStat
	for {
		select {
		case p, ok := <-retChan:
			if !ok {
				if errLogicalDrives != nil {
					return ret, errLogicalDrives
				}
				return ret, warnings.Reference()
			}
			ret = append(ret, p)
		case <-ctx.Done():
			return ret, ctx.Err()
		}
	}
}

func IOCountersWithContext(ctx context.Context, names ...string) (map[string]IOCountersStat, error) {
	// https://github.com/giampaolo/psutil/blob/544e9daa4f66a9f80d7bf6c7886d693ee42f0a13/psutil/arch/windows/disk.c#L83
	drivemap := make(map[string]IOCountersStat, 0)
	var diskPerformance diskPerformance

	lpBuffer := make([]uint16, 254)
	lpBufferLen, err := windows.GetLogicalDriveStrings(uint32(len(lpBuffer)), &lpBuffer[0])
	if err != nil {
		return drivemap, err
	}
	for _, v := range lpBuffer[:lpBufferLen] {
		if 'A' <= v && v <= 'Z' {
			path := string(rune(v)) + ":"
			typepath, _ := windows.UTF16PtrFromString(path)
			typeret := windows.GetDriveType(typepath)
			if typeret == 0 {
				return drivemap, windows.GetLastError()
			}
			if typeret != windows.DRIVE_FIXED {
				continue
			}
			szDevice := fmt.Sprintf(`\\.\%s`, path)
			const IOCTL_DISK_PERFORMANCE = 0x70020
			h, err := windows.CreateFile(syscall.StringToUTF16Ptr(szDevice), 0, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.OPEN_EXISTING, 0, 0)
			if err != nil {
				if err == windows.ERROR_FILE_NOT_FOUND {
					continue
				}
				return drivemap, err
			}
			defer windows.CloseHandle(h)

			var diskPerformanceSize uint32
			err = windows.DeviceIoControl(h, IOCTL_DISK_PERFORMANCE, nil, 0, (*byte)(unsafe.Pointer(&diskPerformance)), uint32(unsafe.Sizeof(diskPerformance)), &diskPerformanceSize, nil)
			if err != nil {
				return drivemap, err
			}
			drivemap[path] = IOCountersStat{
				ReadBytes:  uint64(diskPerformance.BytesRead),
				WriteBytes: uint64(diskPerformance.BytesWritten),
				ReadCount:  uint64(diskPerformance.ReadCount),
				WriteCount: uint64(diskPerformance.WriteCount),
				ReadTime:   uint64(diskPerformance.ReadTime / 10000 / 1000), // convert to ms: https://github.com/giampaolo/psutil/issues/1012
				WriteTime:  uint64(diskPerformance.WriteTime / 10000 / 1000),
				Name:       path,
			}
		}
	}
	return drivemap, nil
}

func SerialNumberWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}

func LabelWithContext(ctx context.Context, name string) (string, error) {
	return "", common.ErrNotImplementedError
}
