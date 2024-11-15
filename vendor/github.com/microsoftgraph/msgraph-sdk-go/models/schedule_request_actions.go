package models
type ScheduleRequestActions int

const (
    ADMINASSIGN_SCHEDULEREQUESTACTIONS ScheduleRequestActions = iota
    ADMINUPDATE_SCHEDULEREQUESTACTIONS
    ADMINREMOVE_SCHEDULEREQUESTACTIONS
    SELFACTIVATE_SCHEDULEREQUESTACTIONS
    SELFDEACTIVATE_SCHEDULEREQUESTACTIONS
    ADMINEXTEND_SCHEDULEREQUESTACTIONS
    ADMINRENEW_SCHEDULEREQUESTACTIONS
    SELFEXTEND_SCHEDULEREQUESTACTIONS
    SELFRENEW_SCHEDULEREQUESTACTIONS
    UNKNOWNFUTUREVALUE_SCHEDULEREQUESTACTIONS
)

func (i ScheduleRequestActions) String() string {
    return []string{"adminAssign", "adminUpdate", "adminRemove", "selfActivate", "selfDeactivate", "adminExtend", "adminRenew", "selfExtend", "selfRenew", "unknownFutureValue"}[i]
}
func ParseScheduleRequestActions(v string) (any, error) {
    result := ADMINASSIGN_SCHEDULEREQUESTACTIONS
    switch v {
        case "adminAssign":
            result = ADMINASSIGN_SCHEDULEREQUESTACTIONS
        case "adminUpdate":
            result = ADMINUPDATE_SCHEDULEREQUESTACTIONS
        case "adminRemove":
            result = ADMINREMOVE_SCHEDULEREQUESTACTIONS
        case "selfActivate":
            result = SELFACTIVATE_SCHEDULEREQUESTACTIONS
        case "selfDeactivate":
            result = SELFDEACTIVATE_SCHEDULEREQUESTACTIONS
        case "adminExtend":
            result = ADMINEXTEND_SCHEDULEREQUESTACTIONS
        case "adminRenew":
            result = ADMINRENEW_SCHEDULEREQUESTACTIONS
        case "selfExtend":
            result = SELFEXTEND_SCHEDULEREQUESTACTIONS
        case "selfRenew":
            result = SELFRENEW_SCHEDULEREQUESTACTIONS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SCHEDULEREQUESTACTIONS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeScheduleRequestActions(values []ScheduleRequestActions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ScheduleRequestActions) isMultiValue() bool {
    return false
}
