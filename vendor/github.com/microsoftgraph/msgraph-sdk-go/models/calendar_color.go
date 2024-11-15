package models
type CalendarColor int

const (
    AUTO_CALENDARCOLOR CalendarColor = iota
    LIGHTBLUE_CALENDARCOLOR
    LIGHTGREEN_CALENDARCOLOR
    LIGHTORANGE_CALENDARCOLOR
    LIGHTGRAY_CALENDARCOLOR
    LIGHTYELLOW_CALENDARCOLOR
    LIGHTTEAL_CALENDARCOLOR
    LIGHTPINK_CALENDARCOLOR
    LIGHTBROWN_CALENDARCOLOR
    LIGHTRED_CALENDARCOLOR
    MAXCOLOR_CALENDARCOLOR
)

func (i CalendarColor) String() string {
    return []string{"auto", "lightBlue", "lightGreen", "lightOrange", "lightGray", "lightYellow", "lightTeal", "lightPink", "lightBrown", "lightRed", "maxColor"}[i]
}
func ParseCalendarColor(v string) (any, error) {
    result := AUTO_CALENDARCOLOR
    switch v {
        case "auto":
            result = AUTO_CALENDARCOLOR
        case "lightBlue":
            result = LIGHTBLUE_CALENDARCOLOR
        case "lightGreen":
            result = LIGHTGREEN_CALENDARCOLOR
        case "lightOrange":
            result = LIGHTORANGE_CALENDARCOLOR
        case "lightGray":
            result = LIGHTGRAY_CALENDARCOLOR
        case "lightYellow":
            result = LIGHTYELLOW_CALENDARCOLOR
        case "lightTeal":
            result = LIGHTTEAL_CALENDARCOLOR
        case "lightPink":
            result = LIGHTPINK_CALENDARCOLOR
        case "lightBrown":
            result = LIGHTBROWN_CALENDARCOLOR
        case "lightRed":
            result = LIGHTRED_CALENDARCOLOR
        case "maxColor":
            result = MAXCOLOR_CALENDARCOLOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCalendarColor(values []CalendarColor) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CalendarColor) isMultiValue() bool {
    return false
}
