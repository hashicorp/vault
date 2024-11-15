package models
type SimulationContentSource int

const (
    UNKNOWN_SIMULATIONCONTENTSOURCE SimulationContentSource = iota
    GLOBAL_SIMULATIONCONTENTSOURCE
    TENANT_SIMULATIONCONTENTSOURCE
    UNKNOWNFUTUREVALUE_SIMULATIONCONTENTSOURCE
)

func (i SimulationContentSource) String() string {
    return []string{"unknown", "global", "tenant", "unknownFutureValue"}[i]
}
func ParseSimulationContentSource(v string) (any, error) {
    result := UNKNOWN_SIMULATIONCONTENTSOURCE
    switch v {
        case "unknown":
            result = UNKNOWN_SIMULATIONCONTENTSOURCE
        case "global":
            result = GLOBAL_SIMULATIONCONTENTSOURCE
        case "tenant":
            result = TENANT_SIMULATIONCONTENTSOURCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIMULATIONCONTENTSOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSimulationContentSource(values []SimulationContentSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SimulationContentSource) isMultiValue() bool {
    return false
}
