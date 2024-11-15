package models
// Scheduled Action Type Enum
type DeviceComplianceActionType int

const (
    // No Action
    NOACTION_DEVICECOMPLIANCEACTIONTYPE DeviceComplianceActionType = iota
    // Send Notification
    NOTIFICATION_DEVICECOMPLIANCEACTIONTYPE
    // Block the device in AAD
    BLOCK_DEVICECOMPLIANCEACTIONTYPE
    // Retire the device
    RETIRE_DEVICECOMPLIANCEACTIONTYPE
    // Wipe the device
    WIPE_DEVICECOMPLIANCEACTIONTYPE
    // Remove Resource Access Profiles from the device
    REMOVERESOURCEACCESSPROFILES_DEVICECOMPLIANCEACTIONTYPE
    // Send push notification to device
    PUSHNOTIFICATION_DEVICECOMPLIANCEACTIONTYPE
)

func (i DeviceComplianceActionType) String() string {
    return []string{"noAction", "notification", "block", "retire", "wipe", "removeResourceAccessProfiles", "pushNotification"}[i]
}
func ParseDeviceComplianceActionType(v string) (any, error) {
    result := NOACTION_DEVICECOMPLIANCEACTIONTYPE
    switch v {
        case "noAction":
            result = NOACTION_DEVICECOMPLIANCEACTIONTYPE
        case "notification":
            result = NOTIFICATION_DEVICECOMPLIANCEACTIONTYPE
        case "block":
            result = BLOCK_DEVICECOMPLIANCEACTIONTYPE
        case "retire":
            result = RETIRE_DEVICECOMPLIANCEACTIONTYPE
        case "wipe":
            result = WIPE_DEVICECOMPLIANCEACTIONTYPE
        case "removeResourceAccessProfiles":
            result = REMOVERESOURCEACCESSPROFILES_DEVICECOMPLIANCEACTIONTYPE
        case "pushNotification":
            result = PUSHNOTIFICATION_DEVICECOMPLIANCEACTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceComplianceActionType(values []DeviceComplianceActionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceComplianceActionType) isMultiValue() bool {
    return false
}
