package models
// Supported platform types for policies.
type PolicyPlatformType int

const (
    // Android.
    ANDROID_POLICYPLATFORMTYPE PolicyPlatformType = iota
    // AndroidForWork.
    ANDROIDFORWORK_POLICYPLATFORMTYPE
    // iOS.
    IOS_POLICYPLATFORMTYPE
    // MacOS.
    MACOS_POLICYPLATFORMTYPE
    // WindowsPhone 8.1.
    WINDOWSPHONE81_POLICYPLATFORMTYPE
    // Windows 8.1 and later
    WINDOWS81ANDLATER_POLICYPLATFORMTYPE
    // Windows 10 and later.
    WINDOWS10ANDLATER_POLICYPLATFORMTYPE
    // All platforms.
    ALL_POLICYPLATFORMTYPE
)

func (i PolicyPlatformType) String() string {
    return []string{"android", "androidForWork", "iOS", "macOS", "windowsPhone81", "windows81AndLater", "windows10AndLater", "all"}[i]
}
func ParsePolicyPlatformType(v string) (any, error) {
    result := ANDROID_POLICYPLATFORMTYPE
    switch v {
        case "android":
            result = ANDROID_POLICYPLATFORMTYPE
        case "androidForWork":
            result = ANDROIDFORWORK_POLICYPLATFORMTYPE
        case "iOS":
            result = IOS_POLICYPLATFORMTYPE
        case "macOS":
            result = MACOS_POLICYPLATFORMTYPE
        case "windowsPhone81":
            result = WINDOWSPHONE81_POLICYPLATFORMTYPE
        case "windows81AndLater":
            result = WINDOWS81ANDLATER_POLICYPLATFORMTYPE
        case "windows10AndLater":
            result = WINDOWS10ANDLATER_POLICYPLATFORMTYPE
        case "all":
            result = ALL_POLICYPLATFORMTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePolicyPlatformType(values []PolicyPlatformType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PolicyPlatformType) isMultiValue() bool {
    return false
}
