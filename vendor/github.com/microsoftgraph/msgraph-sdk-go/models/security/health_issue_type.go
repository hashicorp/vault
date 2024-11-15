package security
type HealthIssueType int

const (
    SENSOR_HEALTHISSUETYPE HealthIssueType = iota
    GLOBAL_HEALTHISSUETYPE
    UNKNOWNFUTUREVALUE_HEALTHISSUETYPE
)

func (i HealthIssueType) String() string {
    return []string{"sensor", "global", "unknownFutureValue"}[i]
}
func ParseHealthIssueType(v string) (any, error) {
    result := SENSOR_HEALTHISSUETYPE
    switch v {
        case "sensor":
            result = SENSOR_HEALTHISSUETYPE
        case "global":
            result = GLOBAL_HEALTHISSUETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_HEALTHISSUETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeHealthIssueType(values []HealthIssueType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i HealthIssueType) isMultiValue() bool {
    return false
}
