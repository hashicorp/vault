package models
type TeamworkUserIdentityType int

const (
    AADUSER_TEAMWORKUSERIDENTITYTYPE TeamworkUserIdentityType = iota
    ONPREMISEAADUSER_TEAMWORKUSERIDENTITYTYPE
    ANONYMOUSGUEST_TEAMWORKUSERIDENTITYTYPE
    FEDERATEDUSER_TEAMWORKUSERIDENTITYTYPE
    PERSONALMICROSOFTACCOUNTUSER_TEAMWORKUSERIDENTITYTYPE
    SKYPEUSER_TEAMWORKUSERIDENTITYTYPE
    PHONEUSER_TEAMWORKUSERIDENTITYTYPE
    UNKNOWNFUTUREVALUE_TEAMWORKUSERIDENTITYTYPE
    EMAILUSER_TEAMWORKUSERIDENTITYTYPE
)

func (i TeamworkUserIdentityType) String() string {
    return []string{"aadUser", "onPremiseAadUser", "anonymousGuest", "federatedUser", "personalMicrosoftAccountUser", "skypeUser", "phoneUser", "unknownFutureValue", "emailUser"}[i]
}
func ParseTeamworkUserIdentityType(v string) (any, error) {
    result := AADUSER_TEAMWORKUSERIDENTITYTYPE
    switch v {
        case "aadUser":
            result = AADUSER_TEAMWORKUSERIDENTITYTYPE
        case "onPremiseAadUser":
            result = ONPREMISEAADUSER_TEAMWORKUSERIDENTITYTYPE
        case "anonymousGuest":
            result = ANONYMOUSGUEST_TEAMWORKUSERIDENTITYTYPE
        case "federatedUser":
            result = FEDERATEDUSER_TEAMWORKUSERIDENTITYTYPE
        case "personalMicrosoftAccountUser":
            result = PERSONALMICROSOFTACCOUNTUSER_TEAMWORKUSERIDENTITYTYPE
        case "skypeUser":
            result = SKYPEUSER_TEAMWORKUSERIDENTITYTYPE
        case "phoneUser":
            result = PHONEUSER_TEAMWORKUSERIDENTITYTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TEAMWORKUSERIDENTITYTYPE
        case "emailUser":
            result = EMAILUSER_TEAMWORKUSERIDENTITYTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTeamworkUserIdentityType(values []TeamworkUserIdentityType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TeamworkUserIdentityType) isMultiValue() bool {
    return false
}
