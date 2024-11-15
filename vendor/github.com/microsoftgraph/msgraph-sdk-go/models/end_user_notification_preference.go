package models
type EndUserNotificationPreference int

const (
    UNKNOWN_ENDUSERNOTIFICATIONPREFERENCE EndUserNotificationPreference = iota
    MICROSOFT_ENDUSERNOTIFICATIONPREFERENCE
    CUSTOM_ENDUSERNOTIFICATIONPREFERENCE
    UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONPREFERENCE
)

func (i EndUserNotificationPreference) String() string {
    return []string{"unknown", "microsoft", "custom", "unknownFutureValue"}[i]
}
func ParseEndUserNotificationPreference(v string) (any, error) {
    result := UNKNOWN_ENDUSERNOTIFICATIONPREFERENCE
    switch v {
        case "unknown":
            result = UNKNOWN_ENDUSERNOTIFICATIONPREFERENCE
        case "microsoft":
            result = MICROSOFT_ENDUSERNOTIFICATIONPREFERENCE
        case "custom":
            result = CUSTOM_ENDUSERNOTIFICATIONPREFERENCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONPREFERENCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEndUserNotificationPreference(values []EndUserNotificationPreference) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EndUserNotificationPreference) isMultiValue() bool {
    return false
}
