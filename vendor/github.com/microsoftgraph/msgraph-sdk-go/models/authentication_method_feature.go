package models
type AuthenticationMethodFeature int

const (
    SSPRREGISTERED_AUTHENTICATIONMETHODFEATURE AuthenticationMethodFeature = iota
    SSPRENABLED_AUTHENTICATIONMETHODFEATURE
    SSPRCAPABLE_AUTHENTICATIONMETHODFEATURE
    PASSWORDLESSCAPABLE_AUTHENTICATIONMETHODFEATURE
    MFACAPABLE_AUTHENTICATIONMETHODFEATURE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODFEATURE
)

func (i AuthenticationMethodFeature) String() string {
    return []string{"ssprRegistered", "ssprEnabled", "ssprCapable", "passwordlessCapable", "mfaCapable", "unknownFutureValue"}[i]
}
func ParseAuthenticationMethodFeature(v string) (any, error) {
    result := SSPRREGISTERED_AUTHENTICATIONMETHODFEATURE
    switch v {
        case "ssprRegistered":
            result = SSPRREGISTERED_AUTHENTICATIONMETHODFEATURE
        case "ssprEnabled":
            result = SSPRENABLED_AUTHENTICATIONMETHODFEATURE
        case "ssprCapable":
            result = SSPRCAPABLE_AUTHENTICATIONMETHODFEATURE
        case "passwordlessCapable":
            result = PASSWORDLESSCAPABLE_AUTHENTICATIONMETHODFEATURE
        case "mfaCapable":
            result = MFACAPABLE_AUTHENTICATIONMETHODFEATURE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODFEATURE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationMethodFeature(values []AuthenticationMethodFeature) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodFeature) isMultiValue() bool {
    return false
}
