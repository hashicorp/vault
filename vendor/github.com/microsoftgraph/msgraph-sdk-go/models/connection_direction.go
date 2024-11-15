package models
type ConnectionDirection int

const (
    UNKNOWN_CONNECTIONDIRECTION ConnectionDirection = iota
    INBOUND_CONNECTIONDIRECTION
    OUTBOUND_CONNECTIONDIRECTION
    UNKNOWNFUTUREVALUE_CONNECTIONDIRECTION
)

func (i ConnectionDirection) String() string {
    return []string{"unknown", "inbound", "outbound", "unknownFutureValue"}[i]
}
func ParseConnectionDirection(v string) (any, error) {
    result := UNKNOWN_CONNECTIONDIRECTION
    switch v {
        case "unknown":
            result = UNKNOWN_CONNECTIONDIRECTION
        case "inbound":
            result = INBOUND_CONNECTIONDIRECTION
        case "outbound":
            result = OUTBOUND_CONNECTIONDIRECTION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONNECTIONDIRECTION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeConnectionDirection(values []ConnectionDirection) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConnectionDirection) isMultiValue() bool {
    return false
}
