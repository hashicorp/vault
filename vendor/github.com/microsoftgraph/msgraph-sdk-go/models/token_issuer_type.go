package models
type TokenIssuerType int

const (
    AZUREAD_TOKENISSUERTYPE TokenIssuerType = iota
    ADFEDERATIONSERVICES_TOKENISSUERTYPE
    UNKNOWNFUTUREVALUE_TOKENISSUERTYPE
    AZUREADBACKUPAUTH_TOKENISSUERTYPE
    ADFEDERATIONSERVICESMFAADAPTER_TOKENISSUERTYPE
    NPSEXTENSION_TOKENISSUERTYPE
)

func (i TokenIssuerType) String() string {
    return []string{"AzureAD", "ADFederationServices", "UnknownFutureValue", "AzureADBackupAuth", "ADFederationServicesMFAAdapter", "NPSExtension"}[i]
}
func ParseTokenIssuerType(v string) (any, error) {
    result := AZUREAD_TOKENISSUERTYPE
    switch v {
        case "AzureAD":
            result = AZUREAD_TOKENISSUERTYPE
        case "ADFederationServices":
            result = ADFEDERATIONSERVICES_TOKENISSUERTYPE
        case "UnknownFutureValue":
            result = UNKNOWNFUTUREVALUE_TOKENISSUERTYPE
        case "AzureADBackupAuth":
            result = AZUREADBACKUPAUTH_TOKENISSUERTYPE
        case "ADFederationServicesMFAAdapter":
            result = ADFEDERATIONSERVICESMFAADAPTER_TOKENISSUERTYPE
        case "NPSExtension":
            result = NPSEXTENSION_TOKENISSUERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeTokenIssuerType(values []TokenIssuerType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i TokenIssuerType) isMultiValue() bool {
    return false
}
