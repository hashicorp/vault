package models
import (
    "math"
    "strings"
)
// Contains properties for Windows device type. Multiple values can be selected. Default value is `none`.
type WindowsDeviceType int

const (
    // No device types supported. Default value.
    NONE_WINDOWSDEVICETYPE = 1
    // Indicates support for Desktop Windows device type.
    DESKTOP_WINDOWSDEVICETYPE = 2
    // Indicates support for Mobile Windows device type.
    MOBILE_WINDOWSDEVICETYPE = 4
    // Indicates support for Holographic Windows device type.
    HOLOGRAPHIC_WINDOWSDEVICETYPE = 8
    // Indicates support for Team Windows device type.
    TEAM_WINDOWSDEVICETYPE = 16
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_WINDOWSDEVICETYPE = 32
)

func (i WindowsDeviceType) String() string {
    var values []string
    options := []string{"none", "desktop", "mobile", "holographic", "team", "unknownFutureValue"}
    for p := 0; p < 6; p++ {
        mantis := WindowsDeviceType(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseWindowsDeviceType(v string) (any, error) {
    var result WindowsDeviceType
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_WINDOWSDEVICETYPE
            case "desktop":
                result |= DESKTOP_WINDOWSDEVICETYPE
            case "mobile":
                result |= MOBILE_WINDOWSDEVICETYPE
            case "holographic":
                result |= HOLOGRAPHIC_WINDOWSDEVICETYPE
            case "team":
                result |= TEAM_WINDOWSDEVICETYPE
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_WINDOWSDEVICETYPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeWindowsDeviceType(values []WindowsDeviceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsDeviceType) isMultiValue() bool {
    return true
}
