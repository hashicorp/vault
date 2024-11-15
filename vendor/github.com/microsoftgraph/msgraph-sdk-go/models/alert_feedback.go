package models
type AlertFeedback int

const (
    UNKNOWN_ALERTFEEDBACK AlertFeedback = iota
    TRUEPOSITIVE_ALERTFEEDBACK
    FALSEPOSITIVE_ALERTFEEDBACK
    BENIGNPOSITIVE_ALERTFEEDBACK
    UNKNOWNFUTUREVALUE_ALERTFEEDBACK
)

func (i AlertFeedback) String() string {
    return []string{"unknown", "truePositive", "falsePositive", "benignPositive", "unknownFutureValue"}[i]
}
func ParseAlertFeedback(v string) (any, error) {
    result := UNKNOWN_ALERTFEEDBACK
    switch v {
        case "unknown":
            result = UNKNOWN_ALERTFEEDBACK
        case "truePositive":
            result = TRUEPOSITIVE_ALERTFEEDBACK
        case "falsePositive":
            result = FALSEPOSITIVE_ALERTFEEDBACK
        case "benignPositive":
            result = BENIGNPOSITIVE_ALERTFEEDBACK
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ALERTFEEDBACK
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAlertFeedback(values []AlertFeedback) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AlertFeedback) isMultiValue() bool {
    return false
}
