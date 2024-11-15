package models
type AppCredentialRestrictionType int

const (
    PASSWORDADDITION_APPCREDENTIALRESTRICTIONTYPE AppCredentialRestrictionType = iota
    PASSWORDLIFETIME_APPCREDENTIALRESTRICTIONTYPE
    SYMMETRICKEYADDITION_APPCREDENTIALRESTRICTIONTYPE
    SYMMETRICKEYLIFETIME_APPCREDENTIALRESTRICTIONTYPE
    CUSTOMPASSWORDADDITION_APPCREDENTIALRESTRICTIONTYPE
    UNKNOWNFUTUREVALUE_APPCREDENTIALRESTRICTIONTYPE
)

func (i AppCredentialRestrictionType) String() string {
    return []string{"passwordAddition", "passwordLifetime", "symmetricKeyAddition", "symmetricKeyLifetime", "customPasswordAddition", "unknownFutureValue"}[i]
}
func ParseAppCredentialRestrictionType(v string) (any, error) {
    result := PASSWORDADDITION_APPCREDENTIALRESTRICTIONTYPE
    switch v {
        case "passwordAddition":
            result = PASSWORDADDITION_APPCREDENTIALRESTRICTIONTYPE
        case "passwordLifetime":
            result = PASSWORDLIFETIME_APPCREDENTIALRESTRICTIONTYPE
        case "symmetricKeyAddition":
            result = SYMMETRICKEYADDITION_APPCREDENTIALRESTRICTIONTYPE
        case "symmetricKeyLifetime":
            result = SYMMETRICKEYLIFETIME_APPCREDENTIALRESTRICTIONTYPE
        case "customPasswordAddition":
            result = CUSTOMPASSWORDADDITION_APPCREDENTIALRESTRICTIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_APPCREDENTIALRESTRICTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAppCredentialRestrictionType(values []AppCredentialRestrictionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AppCredentialRestrictionType) isMultiValue() bool {
    return false
}
