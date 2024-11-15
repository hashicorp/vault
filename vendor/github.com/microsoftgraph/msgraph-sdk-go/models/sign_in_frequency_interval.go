package models
type SignInFrequencyInterval int

const (
    TIMEBASED_SIGNINFREQUENCYINTERVAL SignInFrequencyInterval = iota
    EVERYTIME_SIGNINFREQUENCYINTERVAL
    UNKNOWNFUTUREVALUE_SIGNINFREQUENCYINTERVAL
)

func (i SignInFrequencyInterval) String() string {
    return []string{"timeBased", "everyTime", "unknownFutureValue"}[i]
}
func ParseSignInFrequencyInterval(v string) (any, error) {
    result := TIMEBASED_SIGNINFREQUENCYINTERVAL
    switch v {
        case "timeBased":
            result = TIMEBASED_SIGNINFREQUENCYINTERVAL
        case "everyTime":
            result = EVERYTIME_SIGNINFREQUENCYINTERVAL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIGNINFREQUENCYINTERVAL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSignInFrequencyInterval(values []SignInFrequencyInterval) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SignInFrequencyInterval) isMultiValue() bool {
    return false
}
