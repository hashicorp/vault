package models
// Indicates the type of return code.
type Win32LobAppReturnCodeType int

const (
    // Failed.
    FAILED_WIN32LOBAPPRETURNCODETYPE Win32LobAppReturnCodeType = iota
    // Success.
    SUCCESS_WIN32LOBAPPRETURNCODETYPE
    // Soft-reboot is required.
    SOFTREBOOT_WIN32LOBAPPRETURNCODETYPE
    // Hard-reboot is required.
    HARDREBOOT_WIN32LOBAPPRETURNCODETYPE
    // Retry.
    RETRY_WIN32LOBAPPRETURNCODETYPE
)

func (i Win32LobAppReturnCodeType) String() string {
    return []string{"failed", "success", "softReboot", "hardReboot", "retry"}[i]
}
func ParseWin32LobAppReturnCodeType(v string) (any, error) {
    result := FAILED_WIN32LOBAPPRETURNCODETYPE
    switch v {
        case "failed":
            result = FAILED_WIN32LOBAPPRETURNCODETYPE
        case "success":
            result = SUCCESS_WIN32LOBAPPRETURNCODETYPE
        case "softReboot":
            result = SOFTREBOOT_WIN32LOBAPPRETURNCODETYPE
        case "hardReboot":
            result = HARDREBOOT_WIN32LOBAPPRETURNCODETYPE
        case "retry":
            result = RETRY_WIN32LOBAPPRETURNCODETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWin32LobAppReturnCodeType(values []Win32LobAppReturnCodeType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Win32LobAppReturnCodeType) isMultiValue() bool {
    return false
}
