package models
// Partner state of this tenant.
type MobileThreatPartnerTenantState int

const (
    // Partner is unavailable.
    UNAVAILABLE_MOBILETHREATPARTNERTENANTSTATE MobileThreatPartnerTenantState = iota
    // Partner is available.
    AVAILABLE_MOBILETHREATPARTNERTENANTSTATE
    // Partner is enabled.
    ENABLED_MOBILETHREATPARTNERTENANTSTATE
    // Partner is unresponsive.
    UNRESPONSIVE_MOBILETHREATPARTNERTENANTSTATE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_MOBILETHREATPARTNERTENANTSTATE
)

func (i MobileThreatPartnerTenantState) String() string {
    return []string{"unavailable", "available", "enabled", "unresponsive", "unknownFutureValue"}[i]
}
func ParseMobileThreatPartnerTenantState(v string) (any, error) {
    result := UNAVAILABLE_MOBILETHREATPARTNERTENANTSTATE
    switch v {
        case "unavailable":
            result = UNAVAILABLE_MOBILETHREATPARTNERTENANTSTATE
        case "available":
            result = AVAILABLE_MOBILETHREATPARTNERTENANTSTATE
        case "enabled":
            result = ENABLED_MOBILETHREATPARTNERTENANTSTATE
        case "unresponsive":
            result = UNRESPONSIVE_MOBILETHREATPARTNERTENANTSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MOBILETHREATPARTNERTENANTSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMobileThreatPartnerTenantState(values []MobileThreatPartnerTenantState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i MobileThreatPartnerTenantState) isMultiValue() bool {
    return false
}
