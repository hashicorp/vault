package models
type CalendarSharingAction int

const (
    ACCEPT_CALENDARSHARINGACTION CalendarSharingAction = iota
    ACCEPTANDVIEWCALENDAR_CALENDARSHARINGACTION
    VIEWCALENDAR_CALENDARSHARINGACTION
    ADDTHISCALENDAR_CALENDARSHARINGACTION
)

func (i CalendarSharingAction) String() string {
    return []string{"accept", "acceptAndViewCalendar", "viewCalendar", "addThisCalendar"}[i]
}
func ParseCalendarSharingAction(v string) (any, error) {
    result := ACCEPT_CALENDARSHARINGACTION
    switch v {
        case "accept":
            result = ACCEPT_CALENDARSHARINGACTION
        case "acceptAndViewCalendar":
            result = ACCEPTANDVIEWCALENDAR_CALENDARSHARINGACTION
        case "viewCalendar":
            result = VIEWCALENDAR_CALENDARSHARINGACTION
        case "addThisCalendar":
            result = ADDTHISCALENDAR_CALENDARSHARINGACTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCalendarSharingAction(values []CalendarSharingAction) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CalendarSharingAction) isMultiValue() bool {
    return false
}
