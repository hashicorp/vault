package models
// Indicates the type of restart action.
type Win32LobAppRestartBehavior int

const (
    // Intune will restart the device after running the app installation if the operation returns a reboot code.
    BASEDONRETURNCODE_WIN32LOBAPPRESTARTBEHAVIOR Win32LobAppRestartBehavior = iota
    // Intune will not take any specific action on reboot codes resulting from app installations. Intune will not attempt to suppress restarts for MSI apps.
    ALLOW_WIN32LOBAPPRESTARTBEHAVIOR
    // Intune will attempt to suppress restarts for MSI apps.
    SUPPRESS_WIN32LOBAPPRESTARTBEHAVIOR
    // Intune will force the device to restart immediately after the app installation operation.
    FORCE_WIN32LOBAPPRESTARTBEHAVIOR
)

func (i Win32LobAppRestartBehavior) String() string {
    return []string{"basedOnReturnCode", "allow", "suppress", "force"}[i]
}
func ParseWin32LobAppRestartBehavior(v string) (any, error) {
    result := BASEDONRETURNCODE_WIN32LOBAPPRESTARTBEHAVIOR
    switch v {
        case "basedOnReturnCode":
            result = BASEDONRETURNCODE_WIN32LOBAPPRESTARTBEHAVIOR
        case "allow":
            result = ALLOW_WIN32LOBAPPRESTARTBEHAVIOR
        case "suppress":
            result = SUPPRESS_WIN32LOBAPPRESTARTBEHAVIOR
        case "force":
            result = FORCE_WIN32LOBAPPRESTARTBEHAVIOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppRestartBehavior(values []Win32LobAppRestartBehavior) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppRestartBehavior) isMultiValue() bool {
    return false
}
