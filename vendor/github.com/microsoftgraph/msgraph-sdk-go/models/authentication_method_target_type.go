package models
type AuthenticationMethodTargetType int

const (
    USER_AUTHENTICATIONMETHODTARGETTYPE AuthenticationMethodTargetType = iota
    GROUP_AUTHENTICATIONMETHODTARGETTYPE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODTARGETTYPE
)

func (i AuthenticationMethodTargetType) String() string {
    return []string{"user", "group", "unknownFutureValue"}[i]
}
func ParseAuthenticationMethodTargetType(v string) (any, error) {
    result := USER_AUTHENTICATIONMETHODTARGETTYPE
    switch v {
        case "user":
            result = USER_AUTHENTICATIONMETHODTARGETTYPE
        case "group":
            result = GROUP_AUTHENTICATIONMETHODTARGETTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODTARGETTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationMethodTargetType(values []AuthenticationMethodTargetType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodTargetType) isMultiValue() bool {
    return false
}
