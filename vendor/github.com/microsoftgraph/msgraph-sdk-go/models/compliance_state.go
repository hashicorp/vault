package models
// Compliance state.
type ComplianceState int

const (
    // Unknown.
    UNKNOWN_COMPLIANCESTATE ComplianceState = iota
    // Compliant.
    COMPLIANT_COMPLIANCESTATE
    // Device is non-compliant and is blocked from corporate resources.
    NONCOMPLIANT_COMPLIANCESTATE
    // Conflict with other rules.
    CONFLICT_COMPLIANCESTATE
    // Error.
    ERROR_COMPLIANCESTATE
    // Device is non-compliant but still has access to corporate resources
    INGRACEPERIOD_COMPLIANCESTATE
    // Managed by Config Manager
    CONFIGMANAGER_COMPLIANCESTATE
)

func (i ComplianceState) String() string {
    return []string{"unknown", "compliant", "noncompliant", "conflict", "error", "inGracePeriod", "configManager"}[i]
}
func ParseComplianceState(v string) (any, error) {
    result := UNKNOWN_COMPLIANCESTATE
    switch v {
        case "unknown":
            result = UNKNOWN_COMPLIANCESTATE
        case "compliant":
            result = COMPLIANT_COMPLIANCESTATE
        case "noncompliant":
            result = NONCOMPLIANT_COMPLIANCESTATE
        case "conflict":
            result = CONFLICT_COMPLIANCESTATE
        case "error":
            result = ERROR_COMPLIANCESTATE
        case "inGracePeriod":
            result = INGRACEPERIOD_COMPLIANCESTATE
        case "configManager":
            result = CONFIGMANAGER_COMPLIANCESTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeComplianceState(values []ComplianceState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ComplianceState) isMultiValue() bool {
    return false
}
