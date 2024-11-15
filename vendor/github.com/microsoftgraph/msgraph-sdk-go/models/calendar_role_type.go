package models
type CalendarRoleType int

const (
    NONE_CALENDARROLETYPE CalendarRoleType = iota
    FREEBUSYREAD_CALENDARROLETYPE
    LIMITEDREAD_CALENDARROLETYPE
    READ_CALENDARROLETYPE
    WRITE_CALENDARROLETYPE
    DELEGATEWITHOUTPRIVATEEVENTACCESS_CALENDARROLETYPE
    DELEGATEWITHPRIVATEEVENTACCESS_CALENDARROLETYPE
    CUSTOM_CALENDARROLETYPE
)

func (i CalendarRoleType) String() string {
    return []string{"none", "freeBusyRead", "limitedRead", "read", "write", "delegateWithoutPrivateEventAccess", "delegateWithPrivateEventAccess", "custom"}[i]
}
func ParseCalendarRoleType(v string) (any, error) {
    result := NONE_CALENDARROLETYPE
    switch v {
        case "none":
            result = NONE_CALENDARROLETYPE
        case "freeBusyRead":
            result = FREEBUSYREAD_CALENDARROLETYPE
        case "limitedRead":
            result = LIMITEDREAD_CALENDARROLETYPE
        case "read":
            result = READ_CALENDARROLETYPE
        case "write":
            result = WRITE_CALENDARROLETYPE
        case "delegateWithoutPrivateEventAccess":
            result = DELEGATEWITHOUTPRIVATEEVENTACCESS_CALENDARROLETYPE
        case "delegateWithPrivateEventAccess":
            result = DELEGATEWITHPRIVATEEVENTACCESS_CALENDARROLETYPE
        case "custom":
            result = CUSTOM_CALENDARROLETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCalendarRoleType(values []CalendarRoleType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CalendarRoleType) isMultiValue() bool {
    return false
}
