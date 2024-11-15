package models
type ActivityDomain int

const (
    UNKNOWN_ACTIVITYDOMAIN ActivityDomain = iota
    WORK_ACTIVITYDOMAIN
    PERSONAL_ACTIVITYDOMAIN
    UNRESTRICTED_ACTIVITYDOMAIN
)

func (i ActivityDomain) String() string {
    return []string{"unknown", "work", "personal", "unrestricted"}[i]
}
func ParseActivityDomain(v string) (any, error) {
    result := UNKNOWN_ACTIVITYDOMAIN
    switch v {
        case "unknown":
            result = UNKNOWN_ACTIVITYDOMAIN
        case "work":
            result = WORK_ACTIVITYDOMAIN
        case "personal":
            result = PERSONAL_ACTIVITYDOMAIN
        case "unrestricted":
            result = UNRESTRICTED_ACTIVITYDOMAIN
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeActivityDomain(values []ActivityDomain) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ActivityDomain) isMultiValue() bool {
    return false
}
