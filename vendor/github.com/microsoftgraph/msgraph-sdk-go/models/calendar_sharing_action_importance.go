package models
type CalendarSharingActionImportance int

const (
    PRIMARY_CALENDARSHARINGACTIONIMPORTANCE CalendarSharingActionImportance = iota
    SECONDARY_CALENDARSHARINGACTIONIMPORTANCE
)

func (i CalendarSharingActionImportance) String() string {
    return []string{"primary", "secondary"}[i]
}
func ParseCalendarSharingActionImportance(v string) (any, error) {
    result := PRIMARY_CALENDARSHARINGACTIONIMPORTANCE
    switch v {
        case "primary":
            result = PRIMARY_CALENDARSHARINGACTIONIMPORTANCE
        case "secondary":
            result = SECONDARY_CALENDARSHARINGACTIONIMPORTANCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCalendarSharingActionImportance(values []CalendarSharingActionImportance) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CalendarSharingActionImportance) isMultiValue() bool {
    return false
}
