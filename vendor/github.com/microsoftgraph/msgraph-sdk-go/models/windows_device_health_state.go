package models
import (
    "math"
    "strings"
)
// Computer endpoint protection state
type WindowsDeviceHealthState int

const (
    // Computer is clean and no action is required
    CLEAN_WINDOWSDEVICEHEALTHSTATE = 1
    // Computer is in pending full scan state
    FULLSCANPENDING_WINDOWSDEVICEHEALTHSTATE = 2
    // Computer is in pending reboot state
    REBOOTPENDING_WINDOWSDEVICEHEALTHSTATE = 4
    // Computer is in pending manual steps state
    MANUALSTEPSPENDING_WINDOWSDEVICEHEALTHSTATE = 8
    // Computer is in pending offline scan state
    OFFLINESCANPENDING_WINDOWSDEVICEHEALTHSTATE = 16
    // Computer is in critical failure state
    CRITICAL_WINDOWSDEVICEHEALTHSTATE = 32
)

func (i WindowsDeviceHealthState) String() string {
    var values []string
    options := []string{"clean", "fullScanPending", "rebootPending", "manualStepsPending", "offlineScanPending", "critical"}
    for p := 0; p < 6; p++ {
        mantis := WindowsDeviceHealthState(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsDeviceHealthState(v string) (any, error) {
    var result WindowsDeviceHealthState
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "clean":
                result |= CLEAN_WINDOWSDEVICEHEALTHSTATE
            case "fullScanPending":
                result |= FULLSCANPENDING_WINDOWSDEVICEHEALTHSTATE
            case "rebootPending":
                result |= REBOOTPENDING_WINDOWSDEVICEHEALTHSTATE
            case "manualStepsPending":
                result |= MANUALSTEPSPENDING_WINDOWSDEVICEHEALTHSTATE
            case "offlineScanPending":
                result |= OFFLINESCANPENDING_WINDOWSDEVICEHEALTHSTATE
            case "critical":
                result |= CRITICAL_WINDOWSDEVICEHEALTHSTATE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsDeviceHealthState(values []WindowsDeviceHealthState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsDeviceHealthState) isMultiValue() bool {
    return true
}
