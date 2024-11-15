package externalconnectors
type ConnectionState int

const (
    DRAFT_CONNECTIONSTATE ConnectionState = iota
    READY_CONNECTIONSTATE
    OBSOLETE_CONNECTIONSTATE
    LIMITEXCEEDED_CONNECTIONSTATE
    UNKNOWNFUTUREVALUE_CONNECTIONSTATE
)

func (i ConnectionState) String() string {
    return []string{"draft", "ready", "obsolete", "limitExceeded", "unknownFutureValue"}[i]
}
func ParseConnectionState(v string) (any, error) {
    result := DRAFT_CONNECTIONSTATE
    switch v {
        case "draft":
            result = DRAFT_CONNECTIONSTATE
        case "ready":
            result = READY_CONNECTIONSTATE
        case "obsolete":
            result = OBSOLETE_CONNECTIONSTATE
        case "limitExceeded":
            result = LIMITEXCEEDED_CONNECTIONSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CONNECTIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeConnectionState(values []ConnectionState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConnectionState) isMultiValue() bool {
    return false
}
