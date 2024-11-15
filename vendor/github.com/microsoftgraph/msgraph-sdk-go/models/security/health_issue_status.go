package security
type HealthIssueStatus int

const (
    OPEN_HEALTHISSUESTATUS HealthIssueStatus = iota
    CLOSED_HEALTHISSUESTATUS
    SUPPRESSED_HEALTHISSUESTATUS
    UNKNOWNFUTUREVALUE_HEALTHISSUESTATUS
)

func (i HealthIssueStatus) String() string {
    return []string{"open", "closed", "suppressed", "unknownFutureValue"}[i]
}
func ParseHealthIssueStatus(v string) (any, error) {
    result := OPEN_HEALTHISSUESTATUS
    switch v {
        case "open":
            result = OPEN_HEALTHISSUESTATUS
        case "closed":
            result = CLOSED_HEALTHISSUESTATUS
        case "suppressed":
            result = SUPPRESSED_HEALTHISSUESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HEALTHISSUESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHealthIssueStatus(values []HealthIssueStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HealthIssueStatus) isMultiValue() bool {
    return false
}
