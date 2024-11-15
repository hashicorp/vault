package models
import (
    "math"
    "strings"
)
type AuthenticationMethodModes int

const (
    PASSWORD_AUTHENTICATIONMETHODMODES = 1
    VOICE_AUTHENTICATIONMETHODMODES = 2
    HARDWAREOATH_AUTHENTICATIONMETHODMODES = 4
    SOFTWAREOATH_AUTHENTICATIONMETHODMODES = 8
    SMS_AUTHENTICATIONMETHODMODES = 16
    FIDO2_AUTHENTICATIONMETHODMODES = 32
    WINDOWSHELLOFORBUSINESS_AUTHENTICATIONMETHODMODES = 64
    MICROSOFTAUTHENTICATORPUSH_AUTHENTICATIONMETHODMODES = 128
    DEVICEBASEDPUSH_AUTHENTICATIONMETHODMODES = 256
    TEMPORARYACCESSPASSONETIME_AUTHENTICATIONMETHODMODES = 512
    TEMPORARYACCESSPASSMULTIUSE_AUTHENTICATIONMETHODMODES = 1024
    EMAIL_AUTHENTICATIONMETHODMODES = 2048
    X509CERTIFICATESINGLEFACTOR_AUTHENTICATIONMETHODMODES = 4096
    X509CERTIFICATEMULTIFACTOR_AUTHENTICATIONMETHODMODES = 8192
    FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES = 16384
    FEDERATEDMULTIFACTOR_AUTHENTICATIONMETHODMODES = 32768
    UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODMODES = 65536
)

func (i AuthenticationMethodModes) String() string {
    var values []string
    options := []string{"password", "voice", "hardwareOath", "softwareOath", "sms", "fido2", "windowsHelloForBusiness", "microsoftAuthenticatorPush", "deviceBasedPush", "temporaryAccessPassOneTime", "temporaryAccessPassMultiUse", "email", "x509CertificateSingleFactor", "x509CertificateMultiFactor", "federatedSingleFactor", "federatedMultiFactor", "unknownFutureValue"}
    for p := 0; p < 17; p++ {
        mantis := AuthenticationMethodModes(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseAuthenticationMethodModes(v string) (any, error) {
    var result AuthenticationMethodModes
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "password":
                result |= PASSWORD_AUTHENTICATIONMETHODMODES
            case "voice":
                result |= VOICE_AUTHENTICATIONMETHODMODES
            case "hardwareOath":
                result |= HARDWAREOATH_AUTHENTICATIONMETHODMODES
            case "softwareOath":
                result |= SOFTWAREOATH_AUTHENTICATIONMETHODMODES
            case "sms":
                result |= SMS_AUTHENTICATIONMETHODMODES
            case "fido2":
                result |= FIDO2_AUTHENTICATIONMETHODMODES
            case "windowsHelloForBusiness":
                result |= WINDOWSHELLOFORBUSINESS_AUTHENTICATIONMETHODMODES
            case "microsoftAuthenticatorPush":
                result |= MICROSOFTAUTHENTICATORPUSH_AUTHENTICATIONMETHODMODES
            case "deviceBasedPush":
                result |= DEVICEBASEDPUSH_AUTHENTICATIONMETHODMODES
            case "temporaryAccessPassOneTime":
                result |= TEMPORARYACCESSPASSONETIME_AUTHENTICATIONMETHODMODES
            case "temporaryAccessPassMultiUse":
                result |= TEMPORARYACCESSPASSMULTIUSE_AUTHENTICATIONMETHODMODES
            case "email":
                result |= EMAIL_AUTHENTICATIONMETHODMODES
            case "x509CertificateSingleFactor":
                result |= X509CERTIFICATESINGLEFACTOR_AUTHENTICATIONMETHODMODES
            case "x509CertificateMultiFactor":
                result |= X509CERTIFICATEMULTIFACTOR_AUTHENTICATIONMETHODMODES
            case "federatedSingleFactor":
                result |= FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES
            case "federatedMultiFactor":
                result |= FEDERATEDMULTIFACTOR_AUTHENTICATIONMETHODMODES
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_AUTHENTICATIONMETHODMODES
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeAuthenticationMethodModes(values []AuthenticationMethodModes) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AuthenticationMethodModes) isMultiValue() bool {
    return true
}
