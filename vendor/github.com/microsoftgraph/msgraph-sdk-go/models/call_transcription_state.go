package models
type CallTranscriptionState int

const (
    NOTSTARTED_CALLTRANSCRIPTIONSTATE CallTranscriptionState = iota
    ACTIVE_CALLTRANSCRIPTIONSTATE
    INACTIVE_CALLTRANSCRIPTIONSTATE
    UNKNOWNFUTUREVALUE_CALLTRANSCRIPTIONSTATE
)

func (i CallTranscriptionState) String() string {
    return []string{"notStarted", "active", "inactive", "unknownFutureValue"}[i]
}
func ParseCallTranscriptionState(v string) (any, error) {
    result := NOTSTARTED_CALLTRANSCRIPTIONSTATE
    switch v {
        case "notStarted":
            result = NOTSTARTED_CALLTRANSCRIPTIONSTATE
        case "active":
            result = ACTIVE_CALLTRANSCRIPTIONSTATE
        case "inactive":
            result = INACTIVE_CALLTRANSCRIPTIONSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CALLTRANSCRIPTIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCallTranscriptionState(values []CallTranscriptionState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CallTranscriptionState) isMultiValue() bool {
    return false
}
