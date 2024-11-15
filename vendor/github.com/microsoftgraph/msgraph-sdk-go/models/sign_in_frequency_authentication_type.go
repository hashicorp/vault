package models
type SignInFrequencyAuthenticationType int

const (
    PRIMARYANDSECONDARYAUTHENTICATION_SIGNINFREQUENCYAUTHENTICATIONTYPE SignInFrequencyAuthenticationType = iota
    SECONDARYAUTHENTICATION_SIGNINFREQUENCYAUTHENTICATIONTYPE
    UNKNOWNFUTUREVALUE_SIGNINFREQUENCYAUTHENTICATIONTYPE
)

func (i SignInFrequencyAuthenticationType) String() string {
    return []string{"primaryAndSecondaryAuthentication", "secondaryAuthentication", "unknownFutureValue"}[i]
}
func ParseSignInFrequencyAuthenticationType(v string) (any, error) {
    result := PRIMARYANDSECONDARYAUTHENTICATION_SIGNINFREQUENCYAUTHENTICATIONTYPE
    switch v {
        case "primaryAndSecondaryAuthentication":
            result = PRIMARYANDSECONDARYAUTHENTICATION_SIGNINFREQUENCYAUTHENTICATIONTYPE
        case "secondaryAuthentication":
            result = SECONDARYAUTHENTICATION_SIGNINFREQUENCYAUTHENTICATIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SIGNINFREQUENCYAUTHENTICATIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSignInFrequencyAuthenticationType(values []SignInFrequencyAuthenticationType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SignInFrequencyAuthenticationType) isMultiValue() bool {
    return false
}
