package models
type ScheduleChangeRequestActor int

const (
    SENDER_SCHEDULECHANGEREQUESTACTOR ScheduleChangeRequestActor = iota
    RECIPIENT_SCHEDULECHANGEREQUESTACTOR
    MANAGER_SCHEDULECHANGEREQUESTACTOR
    SYSTEM_SCHEDULECHANGEREQUESTACTOR
    UNKNOWNFUTUREVALUE_SCHEDULECHANGEREQUESTACTOR
)

func (i ScheduleChangeRequestActor) String() string {
    return []string{"sender", "recipient", "manager", "system", "unknownFutureValue"}[i]
}
func ParseScheduleChangeRequestActor(v string) (any, error) {
    result := SENDER_SCHEDULECHANGEREQUESTACTOR
    switch v {
        case "sender":
            result = SENDER_SCHEDULECHANGEREQUESTACTOR
        case "recipient":
            result = RECIPIENT_SCHEDULECHANGEREQUESTACTOR
        case "manager":
            result = MANAGER_SCHEDULECHANGEREQUESTACTOR
        case "system":
            result = SYSTEM_SCHEDULECHANGEREQUESTACTOR
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SCHEDULECHANGEREQUESTACTOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeScheduleChangeRequestActor(values []ScheduleChangeRequestActor) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ScheduleChangeRequestActor) isMultiValue() bool {
    return false
}
