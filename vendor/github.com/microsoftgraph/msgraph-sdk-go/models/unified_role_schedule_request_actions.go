package models
type UnifiedRoleScheduleRequestActions int

const (
    ADMINASSIGN_UNIFIEDROLESCHEDULEREQUESTACTIONS UnifiedRoleScheduleRequestActions = iota
    ADMINUPDATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
    ADMINREMOVE_UNIFIEDROLESCHEDULEREQUESTACTIONS
    SELFACTIVATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
    SELFDEACTIVATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
    ADMINEXTEND_UNIFIEDROLESCHEDULEREQUESTACTIONS
    ADMINRENEW_UNIFIEDROLESCHEDULEREQUESTACTIONS
    SELFEXTEND_UNIFIEDROLESCHEDULEREQUESTACTIONS
    SELFRENEW_UNIFIEDROLESCHEDULEREQUESTACTIONS
    UNKNOWNFUTUREVALUE_UNIFIEDROLESCHEDULEREQUESTACTIONS
)

func (i UnifiedRoleScheduleRequestActions) String() string {
    return []string{"adminAssign", "adminUpdate", "adminRemove", "selfActivate", "selfDeactivate", "adminExtend", "adminRenew", "selfExtend", "selfRenew", "unknownFutureValue"}[i]
}
func ParseUnifiedRoleScheduleRequestActions(v string) (any, error) {
    result := ADMINASSIGN_UNIFIEDROLESCHEDULEREQUESTACTIONS
    switch v {
        case "adminAssign":
            result = ADMINASSIGN_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "adminUpdate":
            result = ADMINUPDATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "adminRemove":
            result = ADMINREMOVE_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "selfActivate":
            result = SELFACTIVATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "selfDeactivate":
            result = SELFDEACTIVATE_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "adminExtend":
            result = ADMINEXTEND_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "adminRenew":
            result = ADMINRENEW_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "selfExtend":
            result = SELFEXTEND_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "selfRenew":
            result = SELFRENEW_UNIFIEDROLESCHEDULEREQUESTACTIONS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_UNIFIEDROLESCHEDULEREQUESTACTIONS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUnifiedRoleScheduleRequestActions(values []UnifiedRoleScheduleRequestActions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UnifiedRoleScheduleRequestActions) isMultiValue() bool {
    return false
}
