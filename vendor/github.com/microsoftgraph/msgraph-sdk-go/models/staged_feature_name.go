package models
type StagedFeatureName int

const (
    PASSTHROUGHAUTHENTICATION_STAGEDFEATURENAME StagedFeatureName = iota
    SEAMLESSSSO_STAGEDFEATURENAME
    PASSWORDHASHSYNC_STAGEDFEATURENAME
    EMAILASALTERNATEID_STAGEDFEATURENAME
    UNKNOWNFUTUREVALUE_STAGEDFEATURENAME
    CERTIFICATEBASEDAUTHENTICATION_STAGEDFEATURENAME
    MULTIFACTORAUTHENTICATION_STAGEDFEATURENAME
)

func (i StagedFeatureName) String() string {
    return []string{"passthroughAuthentication", "seamlessSso", "passwordHashSync", "emailAsAlternateId", "unknownFutureValue", "certificateBasedAuthentication", "multiFactorAuthentication"}[i]
}
func ParseStagedFeatureName(v string) (any, error) {
    result := PASSTHROUGHAUTHENTICATION_STAGEDFEATURENAME
    switch v {
        case "passthroughAuthentication":
            result = PASSTHROUGHAUTHENTICATION_STAGEDFEATURENAME
        case "seamlessSso":
            result = SEAMLESSSSO_STAGEDFEATURENAME
        case "passwordHashSync":
            result = PASSWORDHASHSYNC_STAGEDFEATURENAME
        case "emailAsAlternateId":
            result = EMAILASALTERNATEID_STAGEDFEATURENAME
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_STAGEDFEATURENAME
        case "certificateBasedAuthentication":
            result = CERTIFICATEBASEDAUTHENTICATION_STAGEDFEATURENAME
        case "multiFactorAuthentication":
            result = MULTIFACTORAUTHENTICATION_STAGEDFEATURENAME
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeStagedFeatureName(values []StagedFeatureName) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i StagedFeatureName) isMultiValue() bool {
    return false
}
