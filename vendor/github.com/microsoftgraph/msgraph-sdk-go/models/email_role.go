package models
type EmailRole int

const (
    UNKNOWN_EMAILROLE EmailRole = iota
    SENDER_EMAILROLE
    RECIPIENT_EMAILROLE
    UNKNOWNFUTUREVALUE_EMAILROLE
)

func (i EmailRole) String() string {
    return []string{"unknown", "sender", "recipient", "unknownFutureValue"}[i]
}
func ParseEmailRole(v string) (any, error) {
    result := UNKNOWN_EMAILROLE
    switch v {
        case "unknown":
            result = UNKNOWN_EMAILROLE
        case "sender":
            result = SENDER_EMAILROLE
        case "recipient":
            result = RECIPIENT_EMAILROLE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_EMAILROLE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEmailRole(values []EmailRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EmailRole) isMultiValue() bool {
    return false
}
