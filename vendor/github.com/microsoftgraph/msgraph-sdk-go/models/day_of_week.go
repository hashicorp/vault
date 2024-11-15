package models
type DayOfWeek int

const (
    SUNDAY_DAYOFWEEK DayOfWeek = iota
    MONDAY_DAYOFWEEK
    TUESDAY_DAYOFWEEK
    WEDNESDAY_DAYOFWEEK
    THURSDAY_DAYOFWEEK
    FRIDAY_DAYOFWEEK
    SATURDAY_DAYOFWEEK
)

func (i DayOfWeek) String() string {
    return []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}[i]
}
func ParseDayOfWeek(v string) (any, error) {
    result := SUNDAY_DAYOFWEEK
    switch v {
        case "sunday":
            result = SUNDAY_DAYOFWEEK
        case "monday":
            result = MONDAY_DAYOFWEEK
        case "tuesday":
            result = TUESDAY_DAYOFWEEK
        case "wednesday":
            result = WEDNESDAY_DAYOFWEEK
        case "thursday":
            result = THURSDAY_DAYOFWEEK
        case "friday":
            result = FRIDAY_DAYOFWEEK
        case "saturday":
            result = SATURDAY_DAYOFWEEK
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDayOfWeek(values []DayOfWeek) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DayOfWeek) isMultiValue() bool {
    return false
}
