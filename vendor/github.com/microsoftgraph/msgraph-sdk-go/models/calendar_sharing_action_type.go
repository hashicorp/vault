package models
type CalendarSharingActionType int

const (
    ACCEPT_CALENDARSHARINGACTIONTYPE CalendarSharingActionType = iota
)

func (i CalendarSharingActionType) String() string {
    return []string{"accept"}[i]
}
func ParseCalendarSharingActionType(v string) (any, error) {
    result := ACCEPT_CALENDARSHARINGACTIONTYPE
    switch v {
        case "accept":
            result = ACCEPT_CALENDARSHARINGACTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCalendarSharingActionType(values []CalendarSharingActionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CalendarSharingActionType) isMultiValue() bool {
    return false
}
