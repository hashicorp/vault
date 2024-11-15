package models
type SigninFrequencyType int

const (
    DAYS_SIGNINFREQUENCYTYPE SigninFrequencyType = iota
    HOURS_SIGNINFREQUENCYTYPE
)

func (i SigninFrequencyType) String() string {
    return []string{"days", "hours"}[i]
}
func ParseSigninFrequencyType(v string) (any, error) {
    result := DAYS_SIGNINFREQUENCYTYPE
    switch v {
        case "days":
            result = DAYS_SIGNINFREQUENCYTYPE
        case "hours":
            result = HOURS_SIGNINFREQUENCYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSigninFrequencyType(values []SigninFrequencyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SigninFrequencyType) isMultiValue() bool {
    return false
}
