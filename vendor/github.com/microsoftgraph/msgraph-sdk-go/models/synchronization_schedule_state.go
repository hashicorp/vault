package models
type SynchronizationScheduleState int

const (
    ACTIVE_SYNCHRONIZATIONSCHEDULESTATE SynchronizationScheduleState = iota
    DISABLED_SYNCHRONIZATIONSCHEDULESTATE
    PAUSED_SYNCHRONIZATIONSCHEDULESTATE
)

func (i SynchronizationScheduleState) String() string {
    return []string{"Active", "Disabled", "Paused"}[i]
}
func ParseSynchronizationScheduleState(v string) (any, error) {
    result := ACTIVE_SYNCHRONIZATIONSCHEDULESTATE
    switch v {
        case "Active":
            result = ACTIVE_SYNCHRONIZATIONSCHEDULESTATE
        case "Disabled":
            result = DISABLED_SYNCHRONIZATIONSCHEDULESTATE
        case "Paused":
            result = PAUSED_SYNCHRONIZATIONSCHEDULESTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSynchronizationScheduleState(values []SynchronizationScheduleState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SynchronizationScheduleState) isMultiValue() bool {
    return false
}
