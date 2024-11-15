package models
type FederatedIdpMfaBehavior int

const (
    ACCEPTIFMFADONEBYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR FederatedIdpMfaBehavior = iota
    ENFORCEMFABYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
    REJECTMFABYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
    UNKNOWNFUTUREVALUE_FEDERATEDIDPMFABEHAVIOR
)

func (i FederatedIdpMfaBehavior) String() string {
    return []string{"acceptIfMfaDoneByFederatedIdp", "enforceMfaByFederatedIdp", "rejectMfaByFederatedIdp", "unknownFutureValue"}[i]
}
func ParseFederatedIdpMfaBehavior(v string) (any, error) {
    result := ACCEPTIFMFADONEBYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
    switch v {
        case "acceptIfMfaDoneByFederatedIdp":
            result = ACCEPTIFMFADONEBYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
        case "enforceMfaByFederatedIdp":
            result = ENFORCEMFABYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
        case "rejectMfaByFederatedIdp":
            result = REJECTMFABYFEDERATEDIDP_FEDERATEDIDPMFABEHAVIOR
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_FEDERATEDIDPMFABEHAVIOR
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeFederatedIdpMfaBehavior(values []FederatedIdpMfaBehavior) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i FederatedIdpMfaBehavior) isMultiValue() bool {
    return false
}
