//go:build aix
// +build aix

package host

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/internal/common"
)

// from https://www.ibm.com/docs/en/aix/7.2?topic=files-utmph-file
const (
	user_PROCESS = 7

	hostTemperatureScale = 1000.0 // Not part of the linked file, but kept just in case it becomes relevant
)

func HostIDWithContext(ctx context.Context) (string, error) {
	out, err := invoke.CommandWithContext(ctx, "uname", "-u")
	if err != nil {
		return "", err
	}

	// The command always returns an extra newline, so we make use of Split() to get only the first line
	return strings.Split(string(out[:]), "\n")[0], nil
}

func numProcs(ctx context.Context) (uint64, error) {
	return 0, common.ErrNotImplementedError
}

func BootTimeWithContext(ctx context.Context) (btime uint64, err error) {
	ut, err := UptimeWithContext(ctx)
	if err != nil {
		return 0, err
	}

	if ut <= 0 {
		return 0, errors.New("Uptime was not set, so cannot calculate boot time from it.")
	}

	ut = ut * 60
	return timeSince(ut), nil
}

// This function takes multiple formats of output frmo the uptime
// command and converts the data into minutes.
// Some examples of uptime output that this command handles:
// 11:54AM   up 13 mins,  1 user,  load average: 2.78, 2.62, 1.79
// 12:41PM   up 1 hr,  1 user,  load average: 2.47, 2.85, 2.83
// 07:43PM   up 5 hrs,  1 user,  load average: 3.27, 2.91, 2.72
// 11:18:23  up 83 days, 18:29,  4 users,  load average: 0.16, 0.03, 0.01
func UptimeWithContext(ctx context.Context) (uint64, error) {
	out, err := invoke.CommandWithContext(ctx, "uptime")
	if err != nil {
		return 0, err
	}

	// Convert our uptime to a series of fields we can extract
	ut := strings.Fields(string(out[:]))

	// Convert the second field value to integer
	var days uint64 = 0
	var hours uint64 = 0
	var minutes uint64 = 0
	if ut[3] == "days," {
		days, err = strconv.ParseUint(ut[2], 10, 64)
		if err != nil {
			return 0, err
		}

		// Split field 4 into hours and minutes
		hm := strings.Split(ut[4], ":")
		hours, err = strconv.ParseUint(hm[0], 10, 64)
		if err != nil {
			return 0, err
		}
		minutes, err = strconv.ParseUint(strings.Replace(hm[1], ",", "", -1), 10, 64)
		if err != nil {
			return 0, err
		}
	} else if ut[3] == "hr," || ut[3] == "hrs," {
		hours, err = strconv.ParseUint(ut[2], 10, 64)
		if err != nil {
			return 0, err
		}
	} else if ut[3] == "mins," {
		minutes, err = strconv.ParseUint(ut[2], 10, 64)
		if err != nil {
			return 0, err
		}
	} else if _, err := strconv.ParseInt(ut[3], 10, 64); err == nil && strings.Contains(ut[2], ":") {
		// Split field 2 into hours and minutes
		hm := strings.Split(ut[2], ":")
		hours, err = strconv.ParseUint(hm[0], 10, 64)
		if err != nil {
			return 0, err
		}
		minutes, err = strconv.ParseUint(strings.Replace(hm[1], ",", "", -1), 10, 64)
		if err != nil {
			return 0, err
		}
	}

	// Stack them all together as minutes
	total_time := (days * 24 * 60) + (hours * 60) + minutes

	return total_time, nil
}

// This is a weak implementation due to the limitations on retrieving this data in AIX
func UsersWithContext(ctx context.Context) ([]UserStat, error) {
	var ret []UserStat
	out, err := invoke.CommandWithContext(ctx, "w")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 3 {
		return []UserStat{}, common.ErrNotImplementedError
	}

	hf := strings.Fields(lines[1]) // headers
	for l := 2; l < len(lines); l++ {
		v := strings.Fields(lines[l]) // values
		us := &UserStat{}
		for i, header := range hf {
			// We're done in any of these use cases
			if i >= len(v) || v[0] == "-" {
				break
			}

			if t, err := strconv.ParseFloat(v[i], 64); err == nil {
				switch header {
				case `User`:
					us.User = strconv.FormatFloat(t, 'f', 1, 64)
				case `tty`:
					us.Terminal = strconv.FormatFloat(t, 'f', 1, 64)
				}
			}
		}

		// Valid User data, so append it
		ret = append(ret, *us)
	}

	return ret, nil
}

// Much of this function could be static. However, to be future proofed, I've made it call the OS for the information in all instances.
func PlatformInformationWithContext(ctx context.Context) (platform string, family string, version string, err error) {
	// Set the platform (which should always, and only be, "AIX") from `uname -s`
	out, err := invoke.CommandWithContext(ctx, "uname", "-s")
	if err != nil {
		return "", "", "", err
	}
	platform = strings.TrimRight(string(out[:]), "\n")

	// Set the family
	family = strings.TrimRight(string(out[:]), "\n")

	// Set the version
	out, err = invoke.CommandWithContext(ctx, "oslevel")
	if err != nil {
		return "", "", "", err
	}
	version = strings.TrimRight(string(out[:]), "\n")

	return platform, family, version, nil
}

func KernelVersionWithContext(ctx context.Context) (version string, err error) {
	out, err := invoke.CommandWithContext(ctx, "oslevel", "-s")
	if err != nil {
		return "", err
	}
	version = strings.TrimRight(string(out[:]), "\n")

	return version, nil
}

func KernelArch() (arch string, err error) {
	out, err := invoke.Command("bootinfo", "-y")
	if err != nil {
		return "", err
	}
	arch = strings.TrimRight(string(out[:]), "\n")

	return arch, nil
}

func VirtualizationWithContext(ctx context.Context) (string, string, error) {
	return "", "", common.ErrNotImplementedError
}

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	return nil, common.ErrNotImplementedError
}
