// +build windows

package host

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/internal/common"
	process "github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows"
)

var (
	procGetSystemTimeAsFileTime = common.Modkernel32.NewProc("GetSystemTimeAsFileTime")
	procGetTickCount32          = common.Modkernel32.NewProc("GetTickCount")
	procGetTickCount64          = common.Modkernel32.NewProc("GetTickCount64")
	procGetNativeSystemInfo     = common.Modkernel32.NewProc("GetNativeSystemInfo")
	procRtlGetVersion           = common.ModNt.NewProc("RtlGetVersion")
)

// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/wdm/ns-wdm-_osversioninfoexw
type osVersionInfoExW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        uint8
	wReserved           uint8
}

type systemInfo struct {
	wProcessorArchitecture      uint16
	wReserved                   uint16
	dwPageSize                  uint32
	lpMinimumApplicationAddress uintptr
	lpMaximumApplicationAddress uintptr
	dwActiveProcessorMask       uintptr
	dwNumberOfProcessors        uint32
	dwProcessorType             uint32
	dwAllocationGranularity     uint32
	wProcessorLevel             uint16
	wProcessorRevision          uint16
}

type msAcpi_ThermalZoneTemperature struct {
	Active             bool
	CriticalTripPoint  uint32
	CurrentTemperature uint32
	InstanceName       string
}

func Info() (*InfoStat, error) {
	return InfoWithContext(context.Background())
}

func InfoWithContext(ctx context.Context) (*InfoStat, error) {
	ret := &InfoStat{
		OS: runtime.GOOS,
	}

	{
		hostname, err := os.Hostname()
		if err == nil {
			ret.Hostname = hostname
		}
	}

	{
		platform, family, version, err := PlatformInformationWithContext(ctx)
		if err == nil {
			ret.Platform = platform
			ret.PlatformFamily = family
			ret.PlatformVersion = version
		} else {
			return ret, err
		}
	}

	{
		kernelArch, err := kernelArch()
		if err == nil {
			ret.KernelArch = kernelArch
		}
	}

	{
		boot, err := BootTimeWithContext(ctx)
		if err == nil {
			ret.BootTime = boot
			ret.Uptime, _ = Uptime()
		}
	}

	{
		hostID, err := getMachineGuid()
		if err == nil {
			ret.HostID = hostID
		}
	}

	{
		procs, err := process.PidsWithContext(ctx)
		if err == nil {
			ret.Procs = uint64(len(procs))
		}
	}

	return ret, nil
}

func getMachineGuid() (string, error) {
	// there has been reports of issues on 32bit using golang.org/x/sys/windows/registry, see https://github.com/shirou/gopsutil/pull/312#issuecomment-277422612
	// for rationale of using windows.RegOpenKeyEx/RegQueryValueEx instead of registry.OpenKey/GetStringValue
	var h windows.Handle
	err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Cryptography`), 0, windows.KEY_READ|windows.KEY_WOW64_64KEY, &h)
	if err != nil {
		return "", err
	}
	defer windows.RegCloseKey(h)

	const windowsRegBufLen = 74 // len(`{`) + len(`abcdefgh-1234-456789012-123345456671` * 2) + len(`}`) // 2 == bytes/UTF16
	const uuidLen = 36

	var regBuf [windowsRegBufLen]uint16
	bufLen := uint32(windowsRegBufLen)
	var valType uint32
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`MachineGuid`), nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return "", err
	}

	hostID := windows.UTF16ToString(regBuf[:])
	hostIDLen := len(hostID)
	if hostIDLen != uuidLen {
		return "", fmt.Errorf("HostID incorrect: %q\n", hostID)
	}

	return strings.ToLower(hostID), nil
}

func Uptime() (uint64, error) {
	return UptimeWithContext(context.Background())
}

func UptimeWithContext(ctx context.Context) (uint64, error) {
	procGetTickCount := procGetTickCount64
	err := procGetTickCount64.Find()
	if err != nil {
		procGetTickCount = procGetTickCount32 // handle WinXP, but keep in mind that "the time will wrap around to zero if the system is run continuously for 49.7 days." from MSDN
	}
	r1, _, lastErr := syscall.Syscall(procGetTickCount.Addr(), 0, 0, 0, 0)
	if lastErr != 0 {
		return 0, lastErr
	}
	return uint64((time.Duration(r1) * time.Millisecond).Seconds()), nil
}

func bootTimeFromUptime(up uint64) uint64 {
	return uint64(time.Now().Unix()) - up
}

// cachedBootTime must be accessed via atomic.Load/StoreUint64
var cachedBootTime uint64

func BootTime() (uint64, error) {
	return BootTimeWithContext(context.Background())
}

func BootTimeWithContext(ctx context.Context) (uint64, error) {
	t := atomic.LoadUint64(&cachedBootTime)
	if t != 0 {
		return t, nil
	}
	up, err := Uptime()
	if err != nil {
		return 0, err
	}
	t = bootTimeFromUptime(up)
	atomic.StoreUint64(&cachedBootTime, t)
	return t, nil
}

func PlatformInformation() (platform string, family string, version string, err error) {
	return PlatformInformationWithContext(context.Background())
}

func PlatformInformationWithContext(ctx context.Context) (platform string, family string, version string, err error) {
	// GetVersionEx lies on Windows 8.1 and returns as Windows 8 if we don't declare compatibility in manifest
	// RtlGetVersion bypasses this lying layer and returns the true Windows version
	// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/wdm/nf-wdm-rtlgetversion
	// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/wdm/ns-wdm-_osversioninfoexw
	var osInfo osVersionInfoExW
	osInfo.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osInfo))
	ret, _, err := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&osInfo)))
	if ret != 0 {
		return
	}

	// Platform
	var h windows.Handle // like getMachineGuid(), we query the registry using the raw windows.RegOpenKeyEx/RegQueryValueEx
	err = windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`), 0, windows.KEY_READ|windows.KEY_WOW64_64KEY, &h)
	if err != nil {
		return
	}
	defer windows.RegCloseKey(h)
	var bufLen uint32
	var valType uint32
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`ProductName`), nil, &valType, nil, &bufLen)
	if err != nil {
		return
	}
	regBuf := make([]uint16, bufLen/2+1)
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`ProductName`), nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return
	}
	platform = windows.UTF16ToString(regBuf[:])
	if !strings.HasPrefix(platform, "Microsoft") {
		platform = "Microsoft " + platform
	}
	err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`CSDVersion`), nil, &valType, nil, &bufLen) // append Service Pack number, only on success
	if err == nil {                                                                                       // don't return an error if only the Service Pack retrieval fails
		regBuf = make([]uint16, bufLen/2+1)
		err = windows.RegQueryValueEx(h, windows.StringToUTF16Ptr(`CSDVersion`), nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
		if err == nil {
			platform += " " + windows.UTF16ToString(regBuf[:])
		}
	}

	// PlatformFamily
	switch osInfo.wProductType {
	case 1:
		family = "Standalone Workstation"
	case 2:
		family = "Server (Domain Controller)"
	case 3:
		family = "Server"
	}

	// Platform Version
	version = fmt.Sprintf("%d.%d.%d Build %d", osInfo.dwMajorVersion, osInfo.dwMinorVersion, osInfo.dwBuildNumber, osInfo.dwBuildNumber)

	return platform, family, version, nil
}

func Users() ([]UserStat, error) {
	return UsersWithContext(context.Background())
}

func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	var ret []UserStat

	return ret, nil
}

func SensorsTemperatures() ([]TemperatureStat, error) {
	return SensorsTemperaturesWithContext(context.Background())
}

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	var ret []TemperatureStat
	var dst []msAcpi_ThermalZoneTemperature
	q := wmi.CreateQuery(&dst, "")
	if err := common.WMIQueryWithContext(ctx, q, &dst, nil, "root/wmi"); err != nil {
		return ret, err
	}

	for _, v := range dst {
		ts := TemperatureStat{
			SensorKey:   v.InstanceName,
			Temperature: kelvinToCelsius(v.CurrentTemperature, 2),
		}
		ret = append(ret, ts)
	}

	return ret, nil
}

func kelvinToCelsius(temp uint32, n int) float64 {
	// wmi return temperature Kelvin * 10, so need to divide the result by 10,
	// and then minus 273.15 to get °Celsius.
	t := float64(temp/10) - 273.15
	n10 := math.Pow10(n)
	return math.Trunc((t+0.5/n10)*n10) / n10
}

func Virtualization() (string, string, error) {
	return VirtualizationWithContext(context.Background())
}

func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	return "", "", common.ErrNotImplementedError
}

func KernelVersion() (string, error) {
	return KernelVersionWithContext(context.Background())
}

func KernelVersionWithContext(ctx context.Context) (string, error) {
	_, _, version, err := PlatformInformation()
	return version, err
}

func kernelArch() (string, error) {
	var systemInfo systemInfo
	procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))

	const (
		PROCESSOR_ARCHITECTURE_INTEL = 0
		PROCESSOR_ARCHITECTURE_ARM   = 5
		PROCESSOR_ARCHITECTURE_ARM64 = 12
		PROCESSOR_ARCHITECTURE_IA64  = 6
		PROCESSOR_ARCHITECTURE_AMD64 = 9
	)
	switch systemInfo.wProcessorArchitecture {
	case PROCESSOR_ARCHITECTURE_INTEL:
		if systemInfo.wProcessorLevel < 3 {
			return "i386", nil
		}
		if systemInfo.wProcessorLevel > 6 {
			return "i686", nil
		}
		return fmt.Sprintf("i%d86", systemInfo.wProcessorLevel), nil
	case PROCESSOR_ARCHITECTURE_ARM:
		return "arm", nil
	case PROCESSOR_ARCHITECTURE_ARM64:
		return "aarch64", nil
	case PROCESSOR_ARCHITECTURE_IA64:
		return "ia64", nil
	case PROCESSOR_ARCHITECTURE_AMD64:
		return "x86_64", nil
	}
	return "", nil
}
