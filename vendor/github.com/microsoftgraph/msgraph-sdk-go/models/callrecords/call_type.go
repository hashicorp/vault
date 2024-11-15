package callrecords
type CallType int

const (
    UNKNOWN_CALLTYPE CallType = iota
    GROUPCALL_CALLTYPE
    PEERTOPEER_CALLTYPE
    UNKNOWNFUTUREVALUE_CALLTYPE
)

func (i CallType) String() string {
    return []string{"unknown", "groupCall", "peerToPeer", "unknownFutureValue"}[i]
}
func ParseCallType(v string) (any, error) {
    result := UNKNOWN_CALLTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_CALLTYPE
        case "groupCall":
            result = GROUPCALL_CALLTYPE
        case "peerToPeer":
            result = PEERTOPEER_CALLTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CALLTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCallType(values []CallType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CallType) isMultiValue() bool {
    return false
}
