// +build linux

package host

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/internal/common"
	"golang.org/x/sys/unix"
)

type LSB struct {
	ID          string
	Release     string
	Codename    string
	Description string
}

// from utmp.h
const USER_PROCESS = 7

func Info() (*InfoStat, error) {
	return InfoWithContext(context.Background())
}

func InfoWithContext(ctx context.Context) (*InfoStat, error) {
	ret := &InfoStat{
		OS: runtime.GOOS,
	}

	hostname, err := os.Hostname()
	if err == nil {
		ret.Hostname = hostname
	}

	platform, family, version, err := PlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = version
	}
	kernelVersion, err := KernelVersion()
	if err == nil {
		ret.KernelVersion = kernelVersion
	}

	kernelArch, err := kernelArch()
	if err == nil {
		ret.KernelArch = kernelArch
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

	if numProcs, err := common.NumProcs(); err == nil {
		ret.Procs = numProcs
	}

	sysProductUUID := common.HostSys("class/dmi/id/product_uuid")
	machineID := common.HostEtc("machine-id")
	procSysKernelRandomBootID := common.HostProc("sys/kernel/random/boot_id")
	switch {
	// In order to read this file, needs to be supported by kernel/arch and run as root
	// so having fallback is important
	case common.PathExists(sysProductUUID):
		lines, err := common.ReadLines(sysProductUUID)
		if err == nil && len(lines) > 0 && lines[0] != "" {
			ret.HostID = strings.ToLower(lines[0])
			break
		}
		fallthrough
	// Fallback on GNU Linux systems with systemd, readable by everyone
	case common.PathExists(machineID):
		lines, err := common.ReadLines(machineID)
		if err == nil && len(lines) > 0 && len(lines[0]) == 32 {
			st := lines[0]
			ret.HostID = fmt.Sprintf("%s-%s-%s-%s-%s", st[0:8], st[8:12], st[12:16], st[16:20], st[20:32])
			break
		}
		fallthrough
	// Not stable between reboot, but better than nothing
	default:
		lines, err := common.ReadLines(procSysKernelRandomBootID)
		if err == nil && len(lines) > 0 && lines[0] != "" {
			ret.HostID = strings.ToLower(lines[0])
		}
	}

	return ret, nil
}

// BootTime returns the system boot time expressed in seconds since the epoch.
func BootTime() (uint64, error) {
	return BootTimeWithContext(context.Background())
}

func BootTimeWithContext(ctx context.Context) (uint64, error) {
	return common.BootTimeWithContext(ctx)
}

func uptime(boot uint64) uint64 {
	return uint64(time.Now().Unix()) - boot
}

func Uptime() (uint64, error) {
	return UptimeWithContext(context.Background())
}

func UptimeWithContext(ctx context.Context) (uint64, error) {
	boot, err := BootTime()
	if err != nil {
		return 0, err
	}
	return uptime(boot), nil
}

func Users() ([]UserStat, error) {
	return UsersWithContext(context.Background())
}

func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	utmpfile := common.HostVar("run/utmp")

	file, err := os.Open(utmpfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	count := len(buf) / sizeOfUtmp

	ret := make([]UserStat, 0, count)

	for i := 0; i < count; i++ {
		b := buf[i*sizeOfUtmp : (i+1)*sizeOfUtmp]

		var u utmp
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

func getLSB() (*LSB, error) {
	ret := &LSB{}
	if common.PathExists(common.HostEtc("lsb-release")) {
		contents, err := common.ReadLines(common.HostEtc("lsb-release"))
		if err != nil {
			return ret, err // return empty
		}
		for _, line := range contents {
			field := strings.Split(line, "=")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "DISTRIB_ID":
				ret.ID = field[1]
			case "DISTRIB_RELEASE":
				ret.Release = field[1]
			case "DISTRIB_CODENAME":
				ret.Codename = field[1]
			case "DISTRIB_DESCRIPTION":
				ret.Description = field[1]
			}
		}
	} else if common.PathExists("/usr/bin/lsb_release") {
		lsb_release, err := exec.LookPath("lsb_release")
		if err != nil {
			return ret, err
		}
		out, err := invoke.Command(lsb_release)
		if err != nil {
			return ret, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			field := strings.Split(line, ":")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "Distributor ID":
				ret.ID = field[1]
			case "Release":
				ret.Release = field[1]
			case "Codename":
				ret.Codename = field[1]
			case "Description":
				ret.Description = field[1]
			}
		}

	}

	return ret, nil
}

func PlatformInformation() (platform string, family string, version string, err error) {
	return PlatformInformationWithContext(context.Background())
}

func PlatformInformationWithContext(ctx context.Context) (platform string, family string, version string, err error) {

	lsb, err := getLSB()
	if err != nil {
		lsb = &LSB{}
	}

	if common.PathExists(common.HostEtc("oracle-release")) {
		platform = "oracle"
		contents, err := common.ReadLines(common.HostEtc("oracle-release"))
		if err == nil {
			version = getRedhatishVersion(contents)
		}

	} else if common.PathExists(common.HostEtc("enterprise-release")) {
		platform = "oracle"
		contents, err := common.ReadLines(common.HostEtc("enterprise-release"))
		if err == nil {
			version = getRedhatishVersion(contents)
		}
	} else if common.PathExists(common.HostEtc("slackware-version")) {
		platform = "slackware"
		contents, err := common.ReadLines(common.HostEtc("slackware-version"))
		if err == nil {
			version = getSlackwareVersion(contents)
		}
	} else if common.PathExists(common.HostEtc("debian_version")) {
		if lsb.ID == "Ubuntu" {
			platform = "ubuntu"
			version = lsb.Release
		} else if lsb.ID == "LinuxMint" {
			platform = "linuxmint"
			version = lsb.Release
		} else {
			if common.PathExists("/usr/bin/raspi-config") {
				platform = "raspbian"
			} else {
				platform = "debian"
			}
			contents, err := common.ReadLines(common.HostEtc("debian_version"))
			if err == nil && len(contents) > 0 && contents[0] != "" {
				version = contents[0]
			}
		}
	} else if common.PathExists(common.HostEtc("redhat-release")) {
		contents, err := common.ReadLines(common.HostEtc("redhat-release"))
		if err == nil {
			version = getRedhatishVersion(contents)
			platform = getRedhatishPlatform(contents)
		}
	} else if common.PathExists(common.HostEtc("system-release")) {
		contents, err := common.ReadLines(common.HostEtc("system-release"))
		if err == nil {
			version = getRedhatishVersion(contents)
			platform = getRedhatishPlatform(contents)
		}
	} else if common.PathExists(common.HostEtc("gentoo-release")) {
		platform = "gentoo"
		contents, err := common.ReadLines(common.HostEtc("gentoo-release"))
		if err == nil {
			version = getRedhatishVersion(contents)
		}
	} else if common.PathExists(common.HostEtc("SuSE-release")) {
		contents, err := common.ReadLines(common.HostEtc("SuSE-release"))
		if err == nil {
			version = getSuseVersion(contents)
			platform = getSusePlatform(contents)
		}
		// TODO: slackware detecion
	} else if common.PathExists(common.HostEtc("arch-release")) {
		platform = "arch"
		version = lsb.Release
	} else if common.PathExists(common.HostEtc("alpine-release")) {
		platform = "alpine"
		contents, err := common.ReadLines(common.HostEtc("alpine-release"))
		if err == nil && len(contents) > 0 && contents[0] != "" {
			version = contents[0]
		}
	} else if common.PathExists(common.HostEtc("os-release")) {
		p, v, err := common.GetOSRelease()
		if err == nil {
			platform = p
			version = v
		}
	} else if lsb.ID == "RedHat" {
		platform = "redhat"
		version = lsb.Release
	} else if lsb.ID == "Amazon" {
		platform = "amazon"
		version = lsb.Release
	} else if lsb.ID == "ScientificSL" {
		platform = "scientific"
		version = lsb.Release
	} else if lsb.ID == "XenServer" {
		platform = "xenserver"
		version = lsb.Release
	} else if lsb.ID != "" {
		platform = strings.ToLower(lsb.ID)
		version = lsb.Release
	}

	switch platform {
	case "debian", "ubuntu", "linuxmint", "raspbian":
		family = "debian"
	case "fedora":
		family = "fedora"
	case "oracle", "centos", "redhat", "scientific", "enterpriseenterprise", "amazon", "xenserver", "cloudlinux", "ibm_powerkvm":
		family = "rhel"
	case "suse", "opensuse", "sles":
		family = "suse"
	case "gentoo":
		family = "gentoo"
	case "slackware":
		family = "slackware"
	case "arch":
		family = "arch"
	case "exherbo":
		family = "exherbo"
	case "alpine":
		family = "alpine"
	case "coreos":
		family = "coreos"
	case "solus":
		family = "solus"
	}

	return platform, family, version, nil

}

func KernelVersion() (version string, err error) {
	return KernelVersionWithContext(context.Background())
}

func KernelVersionWithContext(ctx context.Context) (version string, err error) {
	var utsname unix.Utsname
	err = unix.Uname(&utsname)
	if err != nil {
		return "", err
	}
	return string(utsname.Release[:bytes.IndexByte(utsname.Release[:], 0)]), nil
}

func getSlackwareVersion(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))
	c = strings.Replace(c, "slackware ", "", 1)
	return c
}

func getRedhatishVersion(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))

	if strings.Contains(c, "rawhide") {
		return "rawhide"
	}
	if matches := regexp.MustCompile(`release (\d[\d.]*)`).FindStringSubmatch(c); matches != nil {
		return matches[1]
	}
	return ""
}

func getRedhatishPlatform(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))

	if strings.Contains(c, "red hat") {
		return "redhat"
	}
	f := strings.Split(c, " ")

	return f[0]
}

func getSuseVersion(contents []string) string {
	version := ""
	for _, line := range contents {
		if matches := regexp.MustCompile(`VERSION = ([\d.]+)`).FindStringSubmatch(line); matches != nil {
			version = matches[1]
		} else if matches := regexp.MustCompile(`PATCHLEVEL = ([\d]+)`).FindStringSubmatch(line); matches != nil {
			version = version + "." + matches[1]
		}
	}
	return version
}

func getSusePlatform(contents []string) string {
	c := strings.ToLower(strings.Join(contents, ""))
	if strings.Contains(c, "opensuse") {
		return "opensuse"
	}
	return "suse"
}

func Virtualization() (string, string, error) {
	return VirtualizationWithContext(context.Background())
}

func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	return common.VirtualizationWithContext(ctx)
}

func SensorsTemperatures() ([]TemperatureStat, error) {
	return SensorsTemperaturesWithContext(context.Background())
}

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	var temperatures []TemperatureStat
	files, err := filepath.Glob(common.HostSys("/class/hwmon/hwmon*/temp*_*"))
	if err != nil {
		return temperatures, err
	}
	if len(files) == 0 {
		// CentOS has an intermediate /device directory:
		// https://github.com/giampaolo/psutil/issues/971
		files, err = filepath.Glob(common.HostSys("/class/hwmon/hwmon*/device/temp*_*"))
		if err != nil {
			return temperatures, err
		}
	}
	var warns Warnings

	if len(files) == 0 { // handle distributions without hwmon, like raspbian #391, parse legacy thermal_zone files
		files, err = filepath.Glob(common.HostSys("/class/thermal/thermal_zone*/"))
		if err != nil {
			return temperatures, err
		}
		for _, file := range files {
			// Get the name of the temperature you are reading
			name, err := ioutil.ReadFile(filepath.Join(file, "type"))
			if err != nil {
				warns.Add(err)
				continue
			}
			// Get the temperature reading
			current, err := ioutil.ReadFile(filepath.Join(file, "temp"))
			if err != nil {
				warns.Add(err)
				continue
			}
			temperature, err := strconv.ParseInt(strings.TrimSpace(string(current)), 10, 64)
			if err != nil {
				warns.Add(err)
				continue
			}

			temperatures = append(temperatures, TemperatureStat{
				SensorKey:   strings.TrimSpace(string(name)),
				Temperature: float64(temperature) / 1000.0,
			})
		}
		return temperatures, warns.Reference()
	}

	// example directory
	// device/           temp1_crit_alarm  temp2_crit_alarm  temp3_crit_alarm  temp4_crit_alarm  temp5_crit_alarm  temp6_crit_alarm  temp7_crit_alarm
	// name              temp1_input       temp2_input       temp3_input       temp4_input       temp5_input       temp6_input       temp7_input
	// power/            temp1_label       temp2_label       temp3_label       temp4_label       temp5_label       temp6_label       temp7_label
	// subsystem/        temp1_max         temp2_max         temp3_max         temp4_max         temp5_max         temp6_max         temp7_max
	// temp1_crit        temp2_crit        temp3_crit        temp4_crit        temp5_crit        temp6_crit        temp7_crit        uevent
	for _, file := range files {
		filename := strings.Split(filepath.Base(file), "_")
		if filename[1] == "label" {
			// Do not try to read the temperature of the label file
			continue
		}

		// Get the label of the temperature you are reading
		var label string
		c, _ := ioutil.ReadFile(filepath.Join(filepath.Dir(file), filename[0]+"_label"))
		if c != nil {
			//format the label from "Core 0" to "core0_"
			label = fmt.Sprintf("%s_", strings.Join(strings.Split(strings.TrimSpace(strings.ToLower(string(c))), " "), ""))
		}

		// Get the name of the temperature you are reading
		name, err := ioutil.ReadFile(filepath.Join(filepath.Dir(file), "name"))
		if err != nil {
			warns.Add(err)
			continue
		}

		// Get the temperature reading
		current, err := ioutil.ReadFile(file)
		if err != nil {
			warns.Add(err)
			continue
		}
		temperature, err := strconv.ParseFloat(strings.TrimSpace(string(current)), 64)
		if err != nil {
			warns.Add(err)
			continue
		}

		tempName := strings.TrimSpace(strings.ToLower(string(strings.Join(filename[1:], ""))))
		temperatures = append(temperatures, TemperatureStat{
			SensorKey:   fmt.Sprintf("%s_%s%s", strings.TrimSpace(string(name)), label, tempName),
			Temperature: temperature / 1000.0,
		})
	}
	return temperatures, warns.Reference()
}
