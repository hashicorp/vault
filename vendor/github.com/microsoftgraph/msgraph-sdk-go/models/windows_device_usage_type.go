package models
type WindowsDeviceUsageType int

const (
    // Default. Indicates that a device is a single-user device.
    SINGLEUSER_WINDOWSDEVICEUSAGETYPE WindowsDeviceUsageType = iota
    // Indicates that a device is a multi-user device.
    SHARED_WINDOWSDEVICEUSAGETYPE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_WINDOWSDEVICEUSAGETYPE
)

func (i WindowsDeviceUsageType) String() string {
    return []string{"singleUser", "shared", "unknownFutureValue"}[i]
}
func ParseWindowsDeviceUsageType(v string) (any, error) {
    result := SINGLEUSER_WINDOWSDEVICEUSAGETYPE
    switch v {
        case "singleUser":
            result = SINGLEUSER_WINDOWSDEVICEUSAGETYPE
        case "shared":
            result = SHARED_WINDOWSDEVICEUSAGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WINDOWSDEVICEUSAGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsDeviceUsageType(values []WindowsDeviceUsageType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsDeviceUsageType) isMultiValue() bool {
    return false
}
