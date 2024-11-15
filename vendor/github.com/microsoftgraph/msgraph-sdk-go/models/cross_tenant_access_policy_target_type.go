package models
type CrossTenantAccessPolicyTargetType int

const (
    USER_CROSSTENANTACCESSPOLICYTARGETTYPE CrossTenantAccessPolicyTargetType = iota
    GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE
    APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE
    UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETTYPE
)

func (i CrossTenantAccessPolicyTargetType) String() string {
    return []string{"user", "group", "application", "unknownFutureValue"}[i]
}
func ParseCrossTenantAccessPolicyTargetType(v string) (any, error) {
    result := USER_CROSSTENANTACCESSPOLICYTARGETTYPE
    switch v {
        case "user":
            result = USER_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "group":
            result = GROUP_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "application":
            result = APPLICATION_CROSSTENANTACCESSPOLICYTARGETTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CROSSTENANTACCESSPOLICYTARGETTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCrossTenantAccessPolicyTargetType(values []CrossTenantAccessPolicyTargetType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CrossTenantAccessPolicyTargetType) isMultiValue() bool {
    return false
}
