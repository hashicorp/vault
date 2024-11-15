package models
// Possible ways of adding a mobile device to management.
type DeviceEnrollmentType int

const (
    // Default value, enrollment type was not collected.
    UNKNOWN_DEVICEENROLLMENTTYPE DeviceEnrollmentType = iota
    // User driven enrollment through BYOD channel.
    USERENROLLMENT_DEVICEENROLLMENTTYPE
    // User enrollment with a device enrollment manager account.
    DEVICEENROLLMENTMANAGER_DEVICEENROLLMENTTYPE
    // Apple bulk enrollment with user challenge. (DEP, Apple Configurator)
    APPLEBULKWITHUSER_DEVICEENROLLMENTTYPE
    // Apple bulk enrollment without user challenge. (DEP, Apple Configurator, Mobile Config)
    APPLEBULKWITHOUTUSER_DEVICEENROLLMENTTYPE
    // Windows 10 Azure AD Join.
    WINDOWSAZUREADJOIN_DEVICEENROLLMENTTYPE
    // Windows 10 Bulk enrollment through ICD with certificate.
    WINDOWSBULKUSERLESS_DEVICEENROLLMENTTYPE
    // Windows 10 automatic enrollment. (Add work account)
    WINDOWSAUTOENROLLMENT_DEVICEENROLLMENTTYPE
    // Windows 10 bulk Azure AD Join.
    WINDOWSBULKAZUREDOMAINJOIN_DEVICEENROLLMENTTYPE
    // Windows 10 Co-Management triggered by AutoPilot or Group Policy.
    WINDOWSCOMANAGEMENT_DEVICEENROLLMENTTYPE
    // Windows 10 Azure AD Join using Device Auth.
    WINDOWSAZUREADJOINUSINGDEVICEAUTH_DEVICEENROLLMENTTYPE
    // Device managed by Apple user enrollment
    APPLEUSERENROLLMENT_DEVICEENROLLMENTTYPE
    // Device managed by Apple user enrollment with service account
    APPLEUSERENROLLMENTWITHSERVICEACCOUNT_DEVICEENROLLMENTTYPE
)

func (i DeviceEnrollmentType) String() string {
    return []string{"unknown", "userEnrollment", "deviceEnrollmentManager", "appleBulkWithUser", "appleBulkWithoutUser", "windowsAzureADJoin", "windowsBulkUserless", "windowsAutoEnrollment", "windowsBulkAzureDomainJoin", "windowsCoManagement", "windowsAzureADJoinUsingDeviceAuth", "appleUserEnrollment", "appleUserEnrollmentWithServiceAccount"}[i]
}
func ParseDeviceEnrollmentType(v string) (any, error) {
    result := UNKNOWN_DEVICEENROLLMENTTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_DEVICEENROLLMENTTYPE
        case "userEnrollment":
            result = USERENROLLMENT_DEVICEENROLLMENTTYPE
        case "deviceEnrollmentManager":
            result = DEVICEENROLLMENTMANAGER_DEVICEENROLLMENTTYPE
        case "appleBulkWithUser":
            result = APPLEBULKWITHUSER_DEVICEENROLLMENTTYPE
        case "appleBulkWithoutUser":
            result = APPLEBULKWITHOUTUSER_DEVICEENROLLMENTTYPE
        case "windowsAzureADJoin":
            result = WINDOWSAZUREADJOIN_DEVICEENROLLMENTTYPE
        case "windowsBulkUserless":
            result = WINDOWSBULKUSERLESS_DEVICEENROLLMENTTYPE
        case "windowsAutoEnrollment":
            result = WINDOWSAUTOENROLLMENT_DEVICEENROLLMENTTYPE
        case "windowsBulkAzureDomainJoin":
            result = WINDOWSBULKAZUREDOMAINJOIN_DEVICEENROLLMENTTYPE
        case "windowsCoManagement":
            result = WINDOWSCOMANAGEMENT_DEVICEENROLLMENTTYPE
        case "windowsAzureADJoinUsingDeviceAuth":
            result = WINDOWSAZUREADJOINUSINGDEVICEAUTH_DEVICEENROLLMENTTYPE
        case "appleUserEnrollment":
            result = APPLEUSERENROLLMENT_DEVICEENROLLMENTTYPE
        case "appleUserEnrollmentWithServiceAccount":
            result = APPLEUSERENROLLMENTWITHSERVICEACCOUNT_DEVICEENROLLMENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceEnrollmentType(values []DeviceEnrollmentType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceEnrollmentType) isMultiValue() bool {
    return false
}
