package models
type AuthenticationMethodKeyStrength int

const (
    NORMAL_AUTHENTICATIONMETHODKEYSTRENGTH AuthenticationMethodKeyStrength = iota
    WEAK_AUTHENTICATIONMETHODKEYSTRENGTH
    UNKNOWN_AUTHENTICATIONMETHODKEYSTRENGTH
)

func (i AuthenticationMethodKeyStrength) String() string {
    return []string{"normal", "weak", "unknown"}[i]
}
func ParseAuthenticationMethodKeyStrength(v string) (any, error) {
    result := NORMAL_AUTHENTICATIONMETHODKEYSTRENGTH
    switch v {
        case "normal":
            result = NORMAL_AUTHENTICATIONMETHODKEYSTRENGTH
        case "weak":
            result = WEAK_AUTHENTICATIONMETHODKEYSTRENGTH
        case "unknown":
            result = UNKNOWN_AUTHENTICATIONMETHODKEYSTRENGTH
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationMethodKeyStrength(values []AuthenticationMethodKeyStrength) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodKeyStrength) isMultiValue() bool {
    return false
}
