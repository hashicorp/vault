package models
type SimulationStatus int

const (
    UNKNOWN_SIMULATIONSTATUS SimulationStatus = iota
    DRAFT_SIMULATIONSTATUS
    RUNNING_SIMULATIONSTATUS
    SCHEDULED_SIMULATIONSTATUS
    SUCCEEDED_SIMULATIONSTATUS
    FAILED_SIMULATIONSTATUS
    CANCELLED_SIMULATIONSTATUS
    EXCLUDED_SIMULATIONSTATUS
    UNKNOWNFUTUREVALUE_SIMULATIONSTATUS
)

func (i SimulationStatus) String() string {
    return []string{"unknown", "draft", "running", "scheduled", "succeeded", "failed", "cancelled", "excluded", "unknownFutureValue"}[i]
}
func ParseSimulationStatus(v string) (any, error) {
    result := UNKNOWN_SIMULATIONSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_SIMULATIONSTATUS
        case "draft":
            result = DRAFT_SIMULATIONSTATUS
        case "running":
            result = RUNNING_SIMULATIONSTATUS
        case "scheduled":
            result = SCHEDULED_SIMULATIONSTATUS
        case "succeeded":
            result = SUCCEEDED_SIMULATIONSTATUS
        case "failed":
            result = FAILED_SIMULATIONSTATUS
        case "cancelled":
            result = CANCELLED_SIMULATIONSTATUS
        case "excluded":
            result = EXCLUDED_SIMULATIONSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIMULATIONSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSimulationStatus(values []SimulationStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SimulationStatus) isMultiValue() bool {
    return false
}
