package models
type RiskState int

const (
    NONE_RISKSTATE RiskState = iota
    CONFIRMEDSAFE_RISKSTATE
    REMEDIATED_RISKSTATE
    DISMISSED_RISKSTATE
    ATRISK_RISKSTATE
    CONFIRMEDCOMPROMISED_RISKSTATE
    UNKNOWNFUTUREVALUE_RISKSTATE
)

func (i RiskState) String() string {
    return []string{"none", "confirmedSafe", "remediated", "dismissed", "atRisk", "confirmedCompromised", "unknownFutureValue"}[i]
}
func ParseRiskState(v string) (any, error) {
    result := NONE_RISKSTATE
    switch v {
        case "none":
            result = NONE_RISKSTATE
        case "confirmedSafe":
            result = CONFIRMEDSAFE_RISKSTATE
        case "remediated":
            result = REMEDIATED_RISKSTATE
        case "dismissed":
            result = DISMISSED_RISKSTATE
        case "atRisk":
            result = ATRISK_RISKSTATE
        case "confirmedCompromised":
            result = CONFIRMEDCOMPROMISED_RISKSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RISKSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRiskState(values []RiskState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RiskState) isMultiValue() bool {
    return false
}
