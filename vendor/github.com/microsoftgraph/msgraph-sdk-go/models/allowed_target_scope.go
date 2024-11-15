package models
type AllowedTargetScope int

const (
    NOTSPECIFIED_ALLOWEDTARGETSCOPE AllowedTargetScope = iota
    SPECIFICDIRECTORYUSERS_ALLOWEDTARGETSCOPE
    SPECIFICCONNECTEDORGANIZATIONUSERS_ALLOWEDTARGETSCOPE
    SPECIFICDIRECTORYSERVICEPRINCIPALS_ALLOWEDTARGETSCOPE
    ALLMEMBERUSERS_ALLOWEDTARGETSCOPE
    ALLDIRECTORYUSERS_ALLOWEDTARGETSCOPE
    ALLDIRECTORYSERVICEPRINCIPALS_ALLOWEDTARGETSCOPE
    ALLCONFIGUREDCONNECTEDORGANIZATIONUSERS_ALLOWEDTARGETSCOPE
    ALLEXTERNALUSERS_ALLOWEDTARGETSCOPE
    UNKNOWNFUTUREVALUE_ALLOWEDTARGETSCOPE
)

func (i AllowedTargetScope) String() string {
    return []string{"notSpecified", "specificDirectoryUsers", "specificConnectedOrganizationUsers", "specificDirectoryServicePrincipals", "allMemberUsers", "allDirectoryUsers", "allDirectoryServicePrincipals", "allConfiguredConnectedOrganizationUsers", "allExternalUsers", "unknownFutureValue"}[i]
}
func ParseAllowedTargetScope(v string) (any, error) {
    result := NOTSPECIFIED_ALLOWEDTARGETSCOPE
    switch v {
        case "notSpecified":
            result = NOTSPECIFIED_ALLOWEDTARGETSCOPE
        case "specificDirectoryUsers":
            result = SPECIFICDIRECTORYUSERS_ALLOWEDTARGETSCOPE
        case "specificConnectedOrganizationUsers":
            result = SPECIFICCONNECTEDORGANIZATIONUSERS_ALLOWEDTARGETSCOPE
        case "specificDirectoryServicePrincipals":
            result = SPECIFICDIRECTORYSERVICEPRINCIPALS_ALLOWEDTARGETSCOPE
        case "allMemberUsers":
            result = ALLMEMBERUSERS_ALLOWEDTARGETSCOPE
        case "allDirectoryUsers":
            result = ALLDIRECTORYUSERS_ALLOWEDTARGETSCOPE
        case "allDirectoryServicePrincipals":
            result = ALLDIRECTORYSERVICEPRINCIPALS_ALLOWEDTARGETSCOPE
        case "allConfiguredConnectedOrganizationUsers":
            result = ALLCONFIGUREDCONNECTEDORGANIZATIONUSERS_ALLOWEDTARGETSCOPE
        case "allExternalUsers":
            result = ALLEXTERNALUSERS_ALLOWEDTARGETSCOPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ALLOWEDTARGETSCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAllowedTargetScope(values []AllowedTargetScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AllowedTargetScope) isMultiValue() bool {
    return false
}
