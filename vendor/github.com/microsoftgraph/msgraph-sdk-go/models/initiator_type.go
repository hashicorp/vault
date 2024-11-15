package models
type InitiatorType int

const (
    USER_INITIATORTYPE InitiatorType = iota
    APPLICATION_INITIATORTYPE
    SYSTEM_INITIATORTYPE
    UNKNOWNFUTUREVALUE_INITIATORTYPE
)

func (i InitiatorType) String() string {
    return []string{"user", "application", "system", "unknownFutureValue"}[i]
}
func ParseInitiatorType(v string) (any, error) {
    result := USER_INITIATORTYPE
    switch v {
        case "user":
            result = USER_INITIATORTYPE
        case "application":
            result = APPLICATION_INITIATORTYPE
        case "system":
            result = SYSTEM_INITIATORTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_INITIATORTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeInitiatorType(values []InitiatorType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i InitiatorType) isMultiValue() bool {
    return false
}
