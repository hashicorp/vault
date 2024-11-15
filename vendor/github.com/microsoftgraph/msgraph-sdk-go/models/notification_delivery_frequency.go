package models
type NotificationDeliveryFrequency int

const (
    UNKNOWN_NOTIFICATIONDELIVERYFREQUENCY NotificationDeliveryFrequency = iota
    WEEKLY_NOTIFICATIONDELIVERYFREQUENCY
    BIWEEKLY_NOTIFICATIONDELIVERYFREQUENCY
    UNKNOWNFUTUREVALUE_NOTIFICATIONDELIVERYFREQUENCY
)

func (i NotificationDeliveryFrequency) String() string {
    return []string{"unknown", "weekly", "biWeekly", "unknownFutureValue"}[i]
}
func ParseNotificationDeliveryFrequency(v string) (any, error) {
    result := UNKNOWN_NOTIFICATIONDELIVERYFREQUENCY
    switch v {
        case "unknown":
            result = UNKNOWN_NOTIFICATIONDELIVERYFREQUENCY
        case "weekly":
            result = WEEKLY_NOTIFICATIONDELIVERYFREQUENCY
        case "biWeekly":
            result = BIWEEKLY_NOTIFICATIONDELIVERYFREQUENCY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_NOTIFICATIONDELIVERYFREQUENCY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeNotificationDeliveryFrequency(values []NotificationDeliveryFrequency) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i NotificationDeliveryFrequency) isMultiValue() bool {
    return false
}
