package models
type WindowsAutopilotDeviceType int

const (
    // Default. Indicates that the device type  is a Windows PC.
    WINDOWSPC_WINDOWSAUTOPILOTDEVICETYPE WindowsAutopilotDeviceType = iota
    // Indicates that the device type is a HoloLens.
    HOLOLENS_WINDOWSAUTOPILOTDEVICETYPE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_WINDOWSAUTOPILOTDEVICETYPE
)

func (i WindowsAutopilotDeviceType) String() string {
    return []string{"windowsPc", "holoLens", "unknownFutureValue"}[i]
}
func ParseWindowsAutopilotDeviceType(v string) (any, error) {
    result := WINDOWSPC_WINDOWSAUTOPILOTDEVICETYPE
    switch v {
        case "windowsPc":
            result = WINDOWSPC_WINDOWSAUTOPILOTDEVICETYPE
        case "holoLens":
            result = HOLOLENS_WINDOWSAUTOPILOTDEVICETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WINDOWSAUTOPILOTDEVICETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsAutopilotDeviceType(values []WindowsAutopilotDeviceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsAutopilotDeviceType) isMultiValue() bool {
    return false
}
