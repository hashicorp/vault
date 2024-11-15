package models
type CrossTenantAccessPolicyTargetConfigurationAccessType int

const (
    ALLOWED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE CrossTenantAccessPolicyTargetConfigurationAccessType = iota
    BLOCKED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
    UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
)

func (i CrossTenantAccessPolicyTargetConfigurationAccessType) String() string {
    return []string{"allowed", "blocked", "unknownFutureValue"}[i]
}
func ParseCrossTenantAccessPolicyTargetConfigurationAccessType(v string) (any, error) {
    result := ALLOWED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
    switch v {
        case "allowed":
            result = ALLOWED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
        case "blocked":
            result = BLOCKED_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETCONFIGURATIONACCESSTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCrossTenantAccessPolicyTargetConfigurationAccessType(values []CrossTenantAccessPolicyTargetConfigurationAccessType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CrossTenantAccessPolicyTargetConfigurationAccessType) isMultiValue() bool {
    return false
}
