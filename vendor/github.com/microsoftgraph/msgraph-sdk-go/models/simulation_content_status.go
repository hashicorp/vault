package models
type SimulationContentStatus int

const (
    UNKNOWN_SIMULATIONCONTENTSTATUS SimulationContentStatus = iota
    DRAFT_SIMULATIONCONTENTSTATUS
    READY_SIMULATIONCONTENTSTATUS
    ARCHIVE_SIMULATIONCONTENTSTATUS
    DELETE_SIMULATIONCONTENTSTATUS
    UNKNOWNFUTUREVALUE_SIMULATIONCONTENTSTATUS
)

func (i SimulationContentStatus) String() string {
    return []string{"unknown", "draft", "ready", "archive", "delete", "unknownFutureValue"}[i]
}
func ParseSimulationContentStatus(v string) (any, error) {
    result := UNKNOWN_SIMULATIONCONTENTSTATUS
    switch v {
        case "unknown":
            result = UNKNOWN_SIMULATIONCONTENTSTATUS
        case "draft":
            result = DRAFT_SIMULATIONCONTENTSTATUS
        case "ready":
            result = READY_SIMULATIONCONTENTSTATUS
        case "archive":
            result = ARCHIVE_SIMULATIONCONTENTSTATUS
        case "delete":
            result = DELETE_SIMULATIONCONTENTSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIMULATIONCONTENTSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSimulationContentStatus(values []SimulationContentStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SimulationContentStatus) isMultiValue() bool {
    return false
}
