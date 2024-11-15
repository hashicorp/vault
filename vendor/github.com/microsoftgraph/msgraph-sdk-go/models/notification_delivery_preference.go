package models
type NotificationDeliveryPreference int

const (
    UNKNOWN_NOTIFICATIONDELIVERYPREFERENCE NotificationDeliveryPreference = iota
    DELIVERIMMEDIETLY_NOTIFICATIONDELIVERYPREFERENCE
    DELIVERAFTERCAMPAIGNEND_NOTIFICATIONDELIVERYPREFERENCE
    UNKNOWNFUTUREVALUE_NOTIFICATIONDELIVERYPREFERENCE
)

func (i NotificationDeliveryPreference) String() string {
    return []string{"unknown", "deliverImmedietly", "deliverAfterCampaignEnd", "unknownFutureValue"}[i]
}
func ParseNotificationDeliveryPreference(v string) (any, error) {
    result := UNKNOWN_NOTIFICATIONDELIVERYPREFERENCE
    switch v {
        case "unknown":
            result = UNKNOWN_NOTIFICATIONDELIVERYPREFERENCE
        case "deliverImmedietly":
            result = DELIVERIMMEDIETLY_NOTIFICATIONDELIVERYPREFERENCE
        case "deliverAfterCampaignEnd":
            result = DELIVERAFTERCAMPAIGNEND_NOTIFICATIONDELIVERYPREFERENCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_NOTIFICATIONDELIVERYPREFERENCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeNotificationDeliveryPreference(values []NotificationDeliveryPreference) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i NotificationDeliveryPreference) isMultiValue() bool {
    return false
}
