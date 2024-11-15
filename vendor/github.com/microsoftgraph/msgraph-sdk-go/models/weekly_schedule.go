package models
// Possible values for a weekly schedule.
type WeeklySchedule int

const (
    // User Defined, default value, no intent.
    USERDEFINED_WEEKLYSCHEDULE WeeklySchedule = iota
    // Everyday.
    EVERYDAY_WEEKLYSCHEDULE
    // Sunday.
    SUNDAY_WEEKLYSCHEDULE
    // Monday.
    MONDAY_WEEKLYSCHEDULE
    // Tuesday.
    TUESDAY_WEEKLYSCHEDULE
    // Wednesday.
    WEDNESDAY_WEEKLYSCHEDULE
    // Thursday.
    THURSDAY_WEEKLYSCHEDULE
    // Friday.
    FRIDAY_WEEKLYSCHEDULE
    // Saturday.
    SATURDAY_WEEKLYSCHEDULE
)

func (i WeeklySchedule) String() string {
    return []string{"userDefined", "everyday", "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}[i]
}
func ParseWeeklySchedule(v string) (any, error) {
    result := USERDEFINED_WEEKLYSCHEDULE
    switch v {
        case "userDefined":
            result = USERDEFINED_WEEKLYSCHEDULE
        case "everyday":
            result = EVERYDAY_WEEKLYSCHEDULE
        case "sunday":
            result = SUNDAY_WEEKLYSCHEDULE
        case "monday":
            result = MONDAY_WEEKLYSCHEDULE
        case "tuesday":
            result = TUESDAY_WEEKLYSCHEDULE
        case "wednesday":
            result = WEDNESDAY_WEEKLYSCHEDULE
        case "thursday":
            result = THURSDAY_WEEKLYSCHEDULE
        case "friday":
            result = FRIDAY_WEEKLYSCHEDULE
        case "saturday":
            result = SATURDAY_WEEKLYSCHEDULE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWeeklySchedule(values []WeeklySchedule) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i WeeklySchedule) isMultiValue() bool {
    return false
}
