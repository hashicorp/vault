// +build darwin

package host

import (
	"bytes"
	"context"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/shirou/gopsutil/internal/common"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/unix"
)

// from utmpx.h
const USER_PROCESS = 7

func Info() (*InfoStat, error) {
	return InfoWithContext(context.Background())
}

func InfoWithContext(ctx context.Context) (*InfoStat, error) {
	ret := &InfoStat{
		OS:             runtime.GOOS,
		PlatformFamily: "darwin",
	}

	hostname, err := os.Hostname()
	if err == nil {
		ret.Hostname = hostname
	}

	kernelVersion, err := KernelVersionWithContext(ctx)
	if err == nil {
		ret.KernelVersion = kernelVersion
	}

	kernelArch, err := kernelArch()
	if err == nil {
		ret.KernelArch = kernelArch
	}

	platform, family, pver, err := PlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = pver
	}

	system, role, err := Virtualization()
	if err == nil {
		ret.VirtualizationSystem = system
		ret.VirtualizationRole = role
	}

	boot, err := BootTime()
	if err == nil {
		ret.BootTime = boot
		ret.Uptime = uptime(boot)
	}

	procs, err := process.Pids()
	if err == nil {
		ret.Procs = uint64(len(procs))
	}

	uuid, err := unix.Sysctl("kern.uuid")
	if err == nil && uuid != "" {
		ret.HostID = strings.ToLower(uuid)
	}

	return ret, nil
}

// cachedBootTime must be accessed via atomic.Load/StoreUint64
var cachedBootTime uint64

func BootTime() (uint64, error) {
	return BootTimeWithContext(context.Background())
}

func BootTimeWithContext(ctx context.Context) (uint64, error) {
	// https://github.com/AaronO/dashd/blob/222e32ef9f7a1f9bea4a8da2c3627c4cb992f860/probe/probe_darwin.go
	t := atomic.LoadUint64(&cachedBootTime)
	if t != 0 {
		return t, nil
	}
	value, err := unix.Sysctl("kern.boottime")
	if err != nil {
		return 0, err
	}
	bytes := []byte(value[:])
	var boottime uint64
	boottime = uint64(bytes[0]) + uint64(bytes[1])*256 + uint64(bytes[2])*256*256 + uint64(bytes[3])*256*256*256

	atomic.StoreUint64(&cachedBootTime, boottime)

	return boottime, nil
}

func uptime(boot uint64) uint64 {
	return uint64(time.Now().Unix()) - boot
}

func Uptime() (uint64, error) {
	return UptimeWithContext(context.Background())
}

func UptimeWithContext(ctx context.Context) (uint64, error) {
	boot, err := BootTimeWithContext(ctx)
	if err != nil {
		return 0, err
	}
	return uptime(boot), nil
}

func Users() ([]UserStat, error) {
	return UsersWithContext(context.Background())
}

func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	utmpfile := "/var/run/utmpx"
	var ret []UserStat

	file, err := os.Open(utmpfile)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return ret, err
	}

	u := Utmpx{}
	entrySize := int(unsafe.Sizeof(u))
	count := len(buf) / entrySize

	for i := 0; i < count; i++ {
		b := buf[i*entrySize : i*entrySize+entrySize]

		var u Utmpx
		br := bytes.NewReader(b)
		err := binary.Read(br, binary.LittleEndian, &u)
		if err != nil {
			continue
		}
		if u.Type != USER_PROCESS {
			continue
		}
		user := UserStat{
			User:     common.IntToString(u.User[:]),
			Terminal: common.IntToString(u.Line[:]),
			Host:     common.IntToString(u.Host[:]),
			Started:  int(u.Tv.Sec),
		}
		ret = append(ret, user)
	}

	return ret, nil

}

func PlatformInformation() (string, string, string, error) {
	return PlatformInformationWithContext(context.Background())
}

func PlatformInformationWithContext(ctx context.Context) (string, string, string, error) {
	platform := ""
	family := ""
	pver := ""

	sw_vers, err := exec.LookPath("sw_vers")
	if err != nil {
		return "", "", "", err
	}

	p, err := unix.Sysctl("kern.ostype")
	if err == nil {
		platform = strings.ToLower(p)
	}

	out, err := invoke.CommandWithContext(ctx, sw_vers, "-productVersion")
	if err == nil {
		pver = strings.ToLower(strings.TrimSpace(string(out)))
	}

	// check if the macos server version file exists
	_, err = os.Stat("/System/Library/CoreServices/ServerVersion.plist")

	// server file doesn't exist
	if os.IsNotExist(err) {
		family = "Standalone Workstation"
	} else {
		family = "Server"
	}

	return platform, family, pver, nil
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
	version, err := unix.Sysctl("kern.osrelease")
	return strings.ToLower(version), err
}
