package models
// Auto restart required notification dismissal method
type AutoRestartNotificationDismissalMethod int

const (
    // Not configured
    NOTCONFIGURED_AUTORESTARTNOTIFICATIONDISMISSALMETHOD AutoRestartNotificationDismissalMethod = iota
    // Auto dismissal Indicates that the notification is automatically dismissed without user intervention
    AUTOMATIC_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
    // User dismissal. Allows the user to dismiss the notification
    USER_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
    // Evolvable enum member
    UNKNOWNFUTUREVALUE_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
)

func (i AutoRestartNotificationDismissalMethod) String() string {
    return []string{"notConfigured", "automatic", "user", "unknownFutureValue"}[i]
}
func ParseAutoRestartNotificationDismissalMethod(v string) (any, error) {
    result := NOTCONFIGURED_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
    switch v {
        case "notConfigured":
            result = NOTCONFIGURED_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
        case "automatic":
            result = AUTOMATIC_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
        case "user":
            result = USER_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTORESTARTNOTIFICATIONDISMISSALMETHOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAutoRestartNotificationDismissalMethod(values []AutoRestartNotificationDismissalMethod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AutoRestartNotificationDismissalMethod) isMultiValue() bool {
    return false
}
