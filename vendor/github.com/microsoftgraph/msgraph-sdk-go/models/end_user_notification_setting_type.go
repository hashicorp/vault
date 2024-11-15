package models
type EndUserNotificationSettingType int

const (
    UNKNOWN_ENDUSERNOTIFICATIONSETTINGTYPE EndUserNotificationSettingType = iota
    NOTRAINING_ENDUSERNOTIFICATIONSETTINGTYPE
    TRAININGSELECTED_ENDUSERNOTIFICATIONSETTINGTYPE
    NONOTIFICATION_ENDUSERNOTIFICATIONSETTINGTYPE
    UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONSETTINGTYPE
)

func (i EndUserNotificationSettingType) String() string {
    return []string{"unknown", "noTraining", "trainingSelected", "noNotification", "unknownFutureValue"}[i]
}
func ParseEndUserNotificationSettingType(v string) (any, error) {
    result := UNKNOWN_ENDUSERNOTIFICATIONSETTINGTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_ENDUSERNOTIFICATIONSETTINGTYPE
        case "noTraining":
            result = NOTRAINING_ENDUSERNOTIFICATIONSETTINGTYPE
        case "trainingSelected":
            result = TRAININGSELECTED_ENDUSERNOTIFICATIONSETTINGTYPE
        case "noNotification":
            result = NONOTIFICATION_ENDUSERNOTIFICATIONSETTINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONSETTINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEndUserNotificationSettingType(values []EndUserNotificationSettingType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EndUserNotificationSettingType) isMultiValue() bool {
    return false
}
