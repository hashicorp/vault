package security
type HostReputationClassification int

const (
    UNKNOWN_HOSTREPUTATIONCLASSIFICATION HostReputationClassification = iota
    NEUTRAL_HOSTREPUTATIONCLASSIFICATION
    SUSPICIOUS_HOSTREPUTATIONCLASSIFICATION
    MALICIOUS_HOSTREPUTATIONCLASSIFICATION
    UNKNOWNFUTUREVALUE_HOSTREPUTATIONCLASSIFICATION
)

func (i HostReputationClassification) String() string {
    return []string{"unknown", "neutral", "suspicious", "malicious", "unknownFutureValue"}[i]
}
func ParseHostReputationClassification(v string) (any, error) {
    result := UNKNOWN_HOSTREPUTATIONCLASSIFICATION
    switch v {
        case "unknown":
            result = UNKNOWN_HOSTREPUTATIONCLASSIFICATION
        case "neutral":
            result = NEUTRAL_HOSTREPUTATIONCLASSIFICATION
        case "suspicious":
            result = SUSPICIOUS_HOSTREPUTATIONCLASSIFICATION
        case "malicious":
            result = MALICIOUS_HOSTREPUTATIONCLASSIFICATION
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HOSTREPUTATIONCLASSIFICATION
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHostReputationClassification(values []HostReputationClassification) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HostReputationClassification) isMultiValue() bool {
    return false
}
