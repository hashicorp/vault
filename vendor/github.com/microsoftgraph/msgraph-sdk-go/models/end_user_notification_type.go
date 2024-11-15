package models
type EndUserNotificationType int

const (
    UNKNOWN_ENDUSERNOTIFICATIONTYPE EndUserNotificationType = iota
    POSITIVEREINFORCEMENT_ENDUSERNOTIFICATIONTYPE
    NOTRAINING_ENDUSERNOTIFICATIONTYPE
    TRAININGASSIGNMENT_ENDUSERNOTIFICATIONTYPE
    TRAININGREMINDER_ENDUSERNOTIFICATIONTYPE
    UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONTYPE
)

func (i EndUserNotificationType) String() string {
    return []string{"unknown", "positiveReinforcement", "noTraining", "trainingAssignment", "trainingReminder", "unknownFutureValue"}[i]
}
func ParseEndUserNotificationType(v string) (any, error) {
    result := UNKNOWN_ENDUSERNOTIFICATIONTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_ENDUSERNOTIFICATIONTYPE
        case "positiveReinforcement":
            result = POSITIVEREINFORCEMENT_ENDUSERNOTIFICATIONTYPE
        case "noTraining":
            result = NOTRAINING_ENDUSERNOTIFICATIONTYPE
        case "trainingAssignment":
            result = TRAININGASSIGNMENT_ENDUSERNOTIFICATIONTYPE
        case "trainingReminder":
            result = TRAININGREMINDER_ENDUSERNOTIFICATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ENDUSERNOTIFICATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEndUserNotificationType(values []EndUserNotificationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EndUserNotificationType) isMultiValue() bool {
    return false
}
