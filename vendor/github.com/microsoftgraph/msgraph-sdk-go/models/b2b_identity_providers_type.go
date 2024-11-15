package models
type B2bIdentityProvidersType int

const (
    AZUREACTIVEDIRECTORY_B2BIDENTITYPROVIDERSTYPE B2bIdentityProvidersType = iota
    EXTERNALFEDERATION_B2BIDENTITYPROVIDERSTYPE
    SOCIALIDENTITYPROVIDERS_B2BIDENTITYPROVIDERSTYPE
    EMAILONETIMEPASSCODE_B2BIDENTITYPROVIDERSTYPE
    MICROSOFTACCOUNT_B2BIDENTITYPROVIDERSTYPE
    DEFAULTCONFIGUREDIDP_B2BIDENTITYPROVIDERSTYPE
    UNKNOWNFUTUREVALUE_B2BIDENTITYPROVIDERSTYPE
)

func (i B2bIdentityProvidersType) String() string {
    return []string{"azureActiveDirectory", "externalFederation", "socialIdentityProviders", "emailOneTimePasscode", "microsoftAccount", "defaultConfiguredIdp", "unknownFutureValue"}[i]
}
func ParseB2bIdentityProvidersType(v string) (any, error) {
    result := AZUREACTIVEDIRECTORY_B2BIDENTITYPROVIDERSTYPE
    switch v {
        case "azureActiveDirectory":
            result = AZUREACTIVEDIRECTORY_B2BIDENTITYPROVIDERSTYPE
        case "externalFederation":
            result = EXTERNALFEDERATION_B2BIDENTITYPROVIDERSTYPE
        case "socialIdentityProviders":
            result = SOCIALIDENTITYPROVIDERS_B2BIDENTITYPROVIDERSTYPE
        case "emailOneTimePasscode":
            result = EMAILONETIMEPASSCODE_B2BIDENTITYPROVIDERSTYPE
        case "microsoftAccount":
            result = MICROSOFTACCOUNT_B2BIDENTITYPROVIDERSTYPE
        case "defaultConfiguredIdp":
            result = DEFAULTCONFIGUREDIDP_B2BIDENTITYPROVIDERSTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_B2BIDENTITYPROVIDERSTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeB2bIdentityProvidersType(values []B2bIdentityProvidersType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i B2bIdentityProvidersType) isMultiValue() bool {
    return false
}
