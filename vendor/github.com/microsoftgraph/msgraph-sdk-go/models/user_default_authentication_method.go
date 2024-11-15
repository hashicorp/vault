package models
type UserDefaultAuthenticationMethod int

const (
    PUSH_USERDEFAULTAUTHENTICATIONMETHOD UserDefaultAuthenticationMethod = iota
    OATH_USERDEFAULTAUTHENTICATIONMETHOD
    VOICEMOBILE_USERDEFAULTAUTHENTICATIONMETHOD
    VOICEALTERNATEMOBILE_USERDEFAULTAUTHENTICATIONMETHOD
    VOICEOFFICE_USERDEFAULTAUTHENTICATIONMETHOD
    SMS_USERDEFAULTAUTHENTICATIONMETHOD
    NONE_USERDEFAULTAUTHENTICATIONMETHOD
    UNKNOWNFUTUREVALUE_USERDEFAULTAUTHENTICATIONMETHOD
)

func (i UserDefaultAuthenticationMethod) String() string {
    return []string{"push", "oath", "voiceMobile", "voiceAlternateMobile", "voiceOffice", "sms", "none", "unknownFutureValue"}[i]
}
func ParseUserDefaultAuthenticationMethod(v string) (any, error) {
    result := PUSH_USERDEFAULTAUTHENTICATIONMETHOD
    switch v {
        case "push":
            result = PUSH_USERDEFAULTAUTHENTICATIONMETHOD
        case "oath":
            result = OATH_USERDEFAULTAUTHENTICATIONMETHOD
        case "voiceMobile":
            result = VOICEMOBILE_USERDEFAULTAUTHENTICATIONMETHOD
        case "voiceAlternateMobile":
            result = VOICEALTERNATEMOBILE_USERDEFAULTAUTHENTICATIONMETHOD
        case "voiceOffice":
            result = VOICEOFFICE_USERDEFAULTAUTHENTICATIONMETHOD
        case "sms":
            result = SMS_USERDEFAULTAUTHENTICATIONMETHOD
        case "none":
            result = NONE_USERDEFAULTAUTHENTICATIONMETHOD
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_USERDEFAULTAUTHENTICATIONMETHOD
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeUserDefaultAuthenticationMethod(values []UserDefaultAuthenticationMethod) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i UserDefaultAuthenticationMethod) isMultiValue() bool {
    return false
}
