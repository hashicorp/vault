package models
type AuthenticationStrengthPolicyType int

const (
    BUILTIN_AUTHENTICATIONSTRENGTHPOLICYTYPE AuthenticationStrengthPolicyType = iota
    CUSTOM_AUTHENTICATIONSTRENGTHPOLICYTYPE
    UNKNOWNFUTUREVALUE_AUTHENTICATIONSTRENGTHPOLICYTYPE
)

func (i AuthenticationStrengthPolicyType) String() string {
    return []string{"builtIn", "custom", "unknownFutureValue"}[i]
}
func ParseAuthenticationStrengthPolicyType(v string) (any, error) {
    result := BUILTIN_AUTHENTICATIONSTRENGTHPOLICYTYPE
    switch v {
        case "builtIn":
            result = BUILTIN_AUTHENTICATIONSTRENGTHPOLICYTYPE
        case "custom":
            result = CUSTOM_AUTHENTICATIONSTRENGTHPOLICYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONSTRENGTHPOLICYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationStrengthPolicyType(values []AuthenticationStrengthPolicyType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationStrengthPolicyType) isMultiValue() bool {
    return false
}
