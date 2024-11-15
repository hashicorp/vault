package models
// Supported platform types.
type DevicePlatformType int

const (
    // Android.
    ANDROID_DEVICEPLATFORMTYPE DevicePlatformType = iota
    // AndroidForWork.
    ANDROIDFORWORK_DEVICEPLATFORMTYPE
    // iOS.
    IOS_DEVICEPLATFORMTYPE
    // MacOS.
    MACOS_DEVICEPLATFORMTYPE
    // WindowsPhone 8.1.
    WINDOWSPHONE81_DEVICEPLATFORMTYPE
    // Windows 8.1 and later
    WINDOWS81ANDLATER_DEVICEPLATFORMTYPE
    // Windows 10 and later.
    WINDOWS10ANDLATER_DEVICEPLATFORMTYPE
    // Android Work Profile.
    ANDROIDWORKPROFILE_DEVICEPLATFORMTYPE
    // Unknown.
    UNKNOWN_DEVICEPLATFORMTYPE
    // Android AOSP.
    ANDROIDAOSP_DEVICEPLATFORMTYPE
    // Indicates Mobile Application Management (MAM) for android devices.
    ANDROIDMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE
    // Indicates Mobile Application Management (MAM) for iOS devices
    IOSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE
    // Evolvable enumeration sentinel value. Do not use
    UNKNOWNFUTUREVALUE_DEVICEPLATFORMTYPE
)

func (i DevicePlatformType) String() string {
    return []string{"android", "androidForWork", "iOS", "macOS", "windowsPhone81", "windows81AndLater", "windows10AndLater", "androidWorkProfile", "unknown", "androidAOSP", "androidMobileApplicationManagement", "iOSMobileApplicationManagement", "unknownFutureValue"}[i]
}
func ParseDevicePlatformType(v string) (any, error) {
    result := ANDROID_DEVICEPLATFORMTYPE
    switch v {
        case "android":
            result = ANDROID_DEVICEPLATFORMTYPE
        case "androidForWork":
            result = ANDROIDFORWORK_DEVICEPLATFORMTYPE
        case "iOS":
            result = IOS_DEVICEPLATFORMTYPE
        case "macOS":
            result = MACOS_DEVICEPLATFORMTYPE
        case "windowsPhone81":
            result = WINDOWSPHONE81_DEVICEPLATFORMTYPE
        case "windows81AndLater":
            result = WINDOWS81ANDLATER_DEVICEPLATFORMTYPE
        case "windows10AndLater":
            result = WINDOWS10ANDLATER_DEVICEPLATFORMTYPE
        case "androidWorkProfile":
            result = ANDROIDWORKPROFILE_DEVICEPLATFORMTYPE
        case "unknown":
            result = UNKNOWN_DEVICEPLATFORMTYPE
        case "androidAOSP":
            result = ANDROIDAOSP_DEVICEPLATFORMTYPE
        case "androidMobileApplicationManagement":
            result = ANDROIDMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE
        case "iOSMobileApplicationManagement":
            result = IOSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEVICEPLATFORMTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDevicePlatformType(values []DevicePlatformType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DevicePlatformType) isMultiValue() bool {
    return false
}
