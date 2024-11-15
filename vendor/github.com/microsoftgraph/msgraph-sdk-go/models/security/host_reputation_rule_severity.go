package security
type HostReputationRuleSeverity int

const (
    UNKNOWN_HOSTREPUTATIONRULESEVERITY HostReputationRuleSeverity = iota
    LOW_HOSTREPUTATIONRULESEVERITY
    MEDIUM_HOSTREPUTATIONRULESEVERITY
    HIGH_HOSTREPUTATIONRULESEVERITY
    UNKNOWNFUTUREVALUE_HOSTREPUTATIONRULESEVERITY
)

func (i HostReputationRuleSeverity) String() string {
    return []string{"unknown", "low", "medium", "high", "unknownFutureValue"}[i]
}
func ParseHostReputationRuleSeverity(v string) (any, error) {
    result := UNKNOWN_HOSTREPUTATIONRULESEVERITY
    switch v {
        case "unknown":
            result = UNKNOWN_HOSTREPUTATIONRULESEVERITY
        case "low":
            result = LOW_HOSTREPUTATIONRULESEVERITY
        case "medium":
            result = MEDIUM_HOSTREPUTATIONRULESEVERITY
        case "high":
            result = HIGH_HOSTREPUTATIONRULESEVERITY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HOSTREPUTATIONRULESEVERITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHostReputationRuleSeverity(values []HostReputationRuleSeverity) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HostReputationRuleSeverity) isMultiValue() bool {
    return false
}
