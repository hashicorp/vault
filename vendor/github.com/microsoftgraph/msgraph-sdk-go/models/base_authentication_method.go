package models
type BaseAuthenticationMethod int

const (
    PASSWORD_BASEAUTHENTICATIONMETHOD BaseAuthenticationMethod = iota
    VOICE_BASEAUTHENTICATIONMETHOD
    HARDWAREOATH_BASEAUTHENTICATIONMETHOD
    SOFTWAREOATH_BASEAUTHENTICATIONMETHOD
    SMS_BASEAUTHENTICATIONMETHOD
    FIDO2_BASEAUTHENTICATIONMETHOD
    WINDOWSHELLOFORBUSINESS_BASEAUTHENTICATIONMETHOD
    MICROSOFTAUTHENTICATOR_BASEAUTHENTICATIONMETHOD
    TEMPORARYACCESSPASS_BASEAUTHENTICATIONMETHOD
    EMAIL_BASEAUTHENTICATIONMETHOD
    X509CERTIFICATE_BASEAUTHENTICATIONMETHOD
    FEDERATION_BASEAUTHENTICATIONMETHOD
    UNKNOWNFUTUREVALUE_BASEAUTHENTICATIONMETHOD
)

func (i BaseAuthenticationMethod) String() string {
    return []string{"password", "voice", "hardwareOath", "softwareOath", "sms", "fido2", "windowsHelloForBusiness", "microsoftAuthenticator", "temporaryAccessPass", "email", "x509Certificate", "federation", "unknownFutureValue"}[i]
}
func ParseBaseAuthenticationMethod(v string) (any, error) {
    result := PASSWORD_BASEAUTHENTICATIONMETHOD
    switch v {
        case "password":
            result = PASSWORD_BASEAUTHENTICATIONMETHOD
        case "voice":
            result = VOICE_BASEAUTHENTICATIONMETHOD
        case "hardwareOath":
            result = HARDWAREOATH_BASEAUTHENTICATIONMETHOD
        case "softwareOath":
            result = SOFTWAREOATH_BASEAUTHENTICATIONMETHOD
        case "sms":
            result = SMS_BASEAUTHENTICATIONMETHOD
        case "fido2":
            result = FIDO2_BASEAUTHENTICATIONMETHOD
        case "windowsHelloForBusiness":
            result = WINDOWSHELLOFORBUSINESS_BASEAUTHENTICATIONMETHOD
        case "microsoftAuthenticator":
            result = MICROSOFTAUTHENTICATOR_BASEAUTHENTICATIONMETHOD
        case "temporaryAccessPass":
            result = TEMPORARYACCESSPASS_BASEAUTHENTICATIONMETHOD
        case "email":
            result = EMAIL_BASEAUTHENTICATIONMETHOD
        case "x509Certificate":
            result = X509CERTIFICATE_BASEAUTHENTICATIONMETHOD
        case "federation":
            result = FEDERATION_BASEAUTHENTICATIONMETHOD
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_BASEAUTHENTICATIONMETHOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeBaseAuthenticationMethod(values []BaseAuthenticationMethod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i BaseAuthenticationMethod) isMultiValue() bool {
    return false
}
