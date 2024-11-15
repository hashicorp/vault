package callrecords
type FailureStage int

const (
    UNKNOWN_FAILURESTAGE FailureStage = iota
    CALLSETUP_FAILURESTAGE
    MIDCALL_FAILURESTAGE
    UNKNOWNFUTUREVALUE_FAILURESTAGE
)

func (i FailureStage) String() string {
    return []string{"unknown", "callSetup", "midcall", "unknownFutureValue"}[i]
}
func ParseFailureStage(v string) (any, error) {
    result := UNKNOWN_FAILURESTAGE
    switch v {
        case "unknown":
            result = UNKNOWN_FAILURESTAGE
        case "callSetup":
            result = CALLSETUP_FAILURESTAGE
        case "midcall":
            result = MIDCALL_FAILURESTAGE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FAILURESTAGE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFailureStage(values []FailureStage) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FailureStage) isMultiValue() bool {
    return false
}
