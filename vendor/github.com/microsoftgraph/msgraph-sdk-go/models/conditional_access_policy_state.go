package models
type ConditionalAccessPolicyState int

const (
    ENABLED_CONDITIONALACCESSPOLICYSTATE ConditionalAccessPolicyState = iota
    DISABLED_CONDITIONALACCESSPOLICYSTATE
    ENABLEDFORREPORTINGBUTNOTENFORCED_CONDITIONALACCESSPOLICYSTATE
)

func (i ConditionalAccessPolicyState) String() string {
    return []string{"enabled", "disabled", "enabledForReportingButNotEnforced"}[i]
}
func ParseConditionalAccessPolicyState(v string) (any, error) {
    result := ENABLED_CONDITIONALACCESSPOLICYSTATE
    switch v {
        case "enabled":
            result = ENABLED_CONDITIONALACCESSPOLICYSTATE
        case "disabled":
            result = DISABLED_CONDITIONALACCESSPOLICYSTATE
        case "enabledForReportingButNotEnforced":
            result = ENABLEDFORREPORTINGBUTNOTENFORCED_CONDITIONALACCESSPOLICYSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeConditionalAccessPolicyState(values []ConditionalAccessPolicyState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ConditionalAccessPolicyState) isMultiValue() bool {
    return false
}
