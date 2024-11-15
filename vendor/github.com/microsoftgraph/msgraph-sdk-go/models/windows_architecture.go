package models
import (
    "math"
    "strings"
)
// Contains properties for Windows architecture.
type WindowsArchitecture int

const (
    // No flags set.
    NONE_WINDOWSARCHITECTURE = 1
    // Whether or not the X86 Windows architecture type is supported.
    X86_WINDOWSARCHITECTURE = 2
    // Whether or not the X64 Windows architecture type is supported.
    X64_WINDOWSARCHITECTURE = 4
    // Whether or not the Arm Windows architecture type is supported.
    ARM_WINDOWSARCHITECTURE = 8
    // Whether or not the Neutral Windows architecture type is supported.
    NEUTRAL_WINDOWSARCHITECTURE = 16
)

func (i WindowsArchitecture) String() string {
    var values []string
    options := []string{"none", "x86", "x64", "arm", "neutral"}
    for p := 0; p < 5; p++ {
        mantis := WindowsArchitecture(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsArchitecture(v string) (any, error) {
    var result WindowsArchitecture
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_WINDOWSARCHITECTURE
            case "x86":
                result |= X86_WINDOWSARCHITECTURE
            case "x64":
                result |= X64_WINDOWSARCHITECTURE
            case "arm":
                result |= ARM_WINDOWSARCHITECTURE
            case "neutral":
                result |= NEUTRAL_WINDOWSARCHITECTURE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsArchitecture(values []WindowsArchitecture) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsArchitecture) isMultiValue() bool {
    return true
}
