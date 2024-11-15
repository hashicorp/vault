package models
type AuthenticationProtocol int

const (
    WSFED_AUTHENTICATIONPROTOCOL AuthenticationProtocol = iota
    SAML_AUTHENTICATIONPROTOCOL
    UNKNOWNFUTUREVALUE_AUTHENTICATIONPROTOCOL
)

func (i AuthenticationProtocol) String() string {
    return []string{"wsFed", "saml", "unknownFutureValue"}[i]
}
func ParseAuthenticationProtocol(v string) (any, error) {
    result := WSFED_AUTHENTICATIONPROTOCOL
    switch v {
        case "wsFed":
            result = WSFED_AUTHENTICATIONPROTOCOL
        case "saml":
            result = SAML_AUTHENTICATIONPROTOCOL
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUTHENTICATIONPROTOCOL
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAuthenticationProtocol(values []AuthenticationProtocol) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationProtocol) isMultiValue() bool {
    return false
}
