package models
type AuthenticationMethodSignInState int

const (
    NOTSUPPORTED_AUTHENTICATIONMETHODSIGNINSTATE AuthenticationMethodSignInState = iota
    NOTALLOWEDBYPOLICY_AUTHENTICATIONMETHODSIGNINSTATE
    NOTENABLED_AUTHENTICATIONMETHODSIGNINSTATE
    PHONENUMBERNOTUNIQUE_AUTHENTICATIONMETHODSIGNINSTATE
    READY_AUTHENTICATIONMETHODSIGNINSTATE
    NOTCONFIGURED_AUTHENTICATIONMETHODSIGNINSTATE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODSIGNINSTATE
)

func (i AuthenticationMethodSignInState) String() string {
    return []string{"notSupported", "notAllowedByPolicy", "notEnabled", "phoneNumberNotUnique", "ready", "notConfigured", "unknownFutureValue"}[i]
}
func ParseAuthenticationMethodSignInState(v string) (any, error) {
    result := NOTSUPPORTED_AUTHENTICATIONMETHODSIGNINSTATE
    switch v {
        case "notSupported":
            result = NOTSUPPORTED_AUTHENTICATIONMETHODSIGNINSTATE
        case "notAllowedByPolicy":
            result = NOTALLOWEDBYPOLICY_AUTHENTICATIONMETHODSIGNINSTATE
        case "notEnabled":
            result = NOTENABLED_AUTHENTICATIONMETHODSIGNINSTATE
        case "phoneNumberNotUnique":
            result = PHONENUMBERNOTUNIQUE_AUTHENTICATIONMETHODSIGNINSTATE
        case "ready":
            result = READY_AUTHENTICATIONMETHODSIGNINSTATE
        case "notConfigured":
            result = NOTCONFIGURED_AUTHENTICATIONMETHODSIGNINSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODSIGNINSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationMethodSignInState(values []AuthenticationMethodSignInState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodSignInState) isMultiValue() bool {
    return false
}
