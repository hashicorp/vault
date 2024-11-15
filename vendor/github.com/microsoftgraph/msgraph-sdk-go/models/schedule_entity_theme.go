package models
type ScheduleEntityTheme int

const (
    WHITE_SCHEDULEENTITYTHEME ScheduleEntityTheme = iota
    BLUE_SCHEDULEENTITYTHEME
    GREEN_SCHEDULEENTITYTHEME
    PURPLE_SCHEDULEENTITYTHEME
    PINK_SCHEDULEENTITYTHEME
    YELLOW_SCHEDULEENTITYTHEME
    GRAY_SCHEDULEENTITYTHEME
    DARKBLUE_SCHEDULEENTITYTHEME
    DARKGREEN_SCHEDULEENTITYTHEME
    DARKPURPLE_SCHEDULEENTITYTHEME
    DARKPINK_SCHEDULEENTITYTHEME
    DARKYELLOW_SCHEDULEENTITYTHEME
    UNKNOWNFUTUREVALUE_SCHEDULEENTITYTHEME
)

func (i ScheduleEntityTheme) String() string {
    return []string{"white", "blue", "green", "purple", "pink", "yellow", "gray", "darkBlue", "darkGreen", "darkPurple", "darkPink", "darkYellow", "unknownFutureValue"}[i]
}
func ParseScheduleEntityTheme(v string) (any, error) {
    result := WHITE_SCHEDULEENTITYTHEME
    switch v {
        case "white":
            result = WHITE_SCHEDULEENTITYTHEME
        case "blue":
            result = BLUE_SCHEDULEENTITYTHEME
        case "green":
            result = GREEN_SCHEDULEENTITYTHEME
        case "purple":
            result = PURPLE_SCHEDULEENTITYTHEME
        case "pink":
            result = PINK_SCHEDULEENTITYTHEME
        case "yellow":
            result = YELLOW_SCHEDULEENTITYTHEME
        case "gray":
            result = GRAY_SCHEDULEENTITYTHEME
        case "darkBlue":
            result = DARKBLUE_SCHEDULEENTITYTHEME
        case "darkGreen":
            result = DARKGREEN_SCHEDULEENTITYTHEME
        case "darkPurple":
            result = DARKPURPLE_SCHEDULEENTITYTHEME
        case "darkPink":
            result = DARKPINK_SCHEDULEENTITYTHEME
        case "darkYellow":
            result = DARKYELLOW_SCHEDULEENTITYTHEME
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SCHEDULEENTITYTHEME
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeScheduleEntityTheme(values []ScheduleEntityTheme) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ScheduleEntityTheme) isMultiValue() bool {
    return false
}
