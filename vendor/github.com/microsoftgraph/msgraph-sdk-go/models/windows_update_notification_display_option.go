package models
// Windows Update Notification Display Options
type WindowsUpdateNotificationDisplayOption int

const (
    // Not configured
    NOTCONFIGURED_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION WindowsUpdateNotificationDisplayOption = iota
    // Use the default Windows Update notifications.
    DEFAULTNOTIFICATIONS_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
    // Turn off all notifications, excluding restart warnings.
    RESTARTWARNINGSONLY_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
    // Turn off all notifications, including restart warnings.
    DISABLEALLNOTIFICATIONS_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
    // Evolvable enum member
    UNKNOWNFUTUREVALUE_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
)

func (i WindowsUpdateNotificationDisplayOption) String() string {
    return []string{"notConfigured", "defaultNotifications", "restartWarningsOnly", "disableAllNotifications", "unknownFutureValue"}[i]
}
func ParseWindowsUpdateNotificationDisplayOption(v string) (any, error) {
    result := NOTCONFIGURED_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
        case "defaultNotifications":
            result = DEFAULTNOTIFICATIONS_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
        case "restartWarningsOnly":
            result = RESTARTWARNINGSONLY_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
        case "disableAllNotifications":
            result = DISABLEALLNOTIFICATIONS_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_WINDOWSUPDATENOTIFICATIONDISPLAYOPTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindowsUpdateNotificationDisplayOption(values []WindowsUpdateNotificationDisplayOption) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WindowsUpdateNotificationDisplayOption) isMultiValue() bool {
    return false
}
