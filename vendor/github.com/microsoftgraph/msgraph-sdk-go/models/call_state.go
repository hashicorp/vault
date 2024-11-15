package models
type CallState int

const (
    INCOMING_CALLSTATE CallState = iota
    ESTABLISHING_CALLSTATE
    ESTABLISHED_CALLSTATE
    HOLD_CALLSTATE
    TRANSFERRING_CALLSTATE
    TRANSFERACCEPTED_CALLSTATE
    REDIRECTING_CALLSTATE
    TERMINATING_CALLSTATE
    TERMINATED_CALLSTATE
    UNKNOWNFUTUREVALUE_CALLSTATE
)

func (i CallState) String() string {
    return []string{"incoming", "establishing", "established", "hold", "transferring", "transferAccepted", "redirecting", "terminating", "terminated", "unknownFutureValue"}[i]
}
func ParseCallState(v string) (any, error) {
    result := INCOMING_CALLSTATE
    switch v {
        case "incoming":
            result = INCOMING_CALLSTATE
        case "establishing":
            result = ESTABLISHING_CALLSTATE
        case "established":
            result = ESTABLISHED_CALLSTATE
        case "hold":
            result = HOLD_CALLSTATE
        case "transferring":
            result = TRANSFERRING_CALLSTATE
        case "transferAccepted":
            result = TRANSFERACCEPTED_CALLSTATE
        case "redirecting":
            result = REDIRECTING_CALLSTATE
        case "terminating":
            result = TERMINATING_CALLSTATE
        case "terminated":
            result = TERMINATED_CALLSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CALLSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCallState(values []CallState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CallState) isMultiValue() bool {
    return false
}
