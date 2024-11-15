package models
// Which branch devices will receive their updates from
type WindowsUpdateType int

const (
    // Allow the user to set.
    USERDEFINED_WINDOWSUPDATETYPE WindowsUpdateType = iota
    // Semi-annual Channel (Targeted). Device gets all applicable feature updates from Semi-annual Channel (Targeted).
    ALL_WINDOWSUPDATETYPE
    // Semi-annual Channel. Device gets feature updates from Semi-annual Channel.
    BUSINESSREADYONLY_WINDOWSUPDATETYPE
    // Windows Insider build - Fast
    WINDOWSINSIDERBUILDFAST_WINDOWSUPDATETYPE
    // Windows Insider build - Slow
    WINDOWSINSIDERBUILDSLOW_WINDOWSUPDATETYPE
    // Release Windows Insider build
    WINDOWSINSIDERBUILDRELEASE_WINDOWSUPDATETYPE
)

func (i WindowsUpdateType) String() string {
    return []string{"userDefined", "all", "businessReadyOnly", "windowsInsiderBuildFast", "windowsInsiderBuildSlow", "windowsInsiderBuildRelease"}[i]
}
func ParseWindowsUpdateType(v string) (any, error) {
    result := USERDEFINED_WINDOWSUPDATETYPE
    switch v {
        case "userDefined":
            result = USERDEFINED_WINDOWSUPDATETYPE
        case "all":
            result = ALL_WINDOWSUPDATETYPE
        case "businessReadyOnly":
            result = BUSINESSREADYONLY_WINDOWSUPDATETYPE
        case "windowsInsiderBuildFast":
            result = WINDOWSINSIDERBUILDFAST_WINDOWSUPDATETYPE
        case "windowsInsiderBuildSlow":
            result = WINDOWSINSIDERBUILDSLOW_WINDOWSUPDATETYPE
        case "windowsInsiderBuildRelease":
            result = WINDOWSINSIDERBUILDRELEASE_WINDOWSUPDATETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsUpdateType(values []WindowsUpdateType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsUpdateType) isMultiValue() bool {
    return false
}
