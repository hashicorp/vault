package security
type ServiceSource int

const (
    UNKNOWN_SERVICESOURCE ServiceSource = iota
    MICROSOFTDEFENDERFORENDPOINT_SERVICESOURCE
    MICROSOFTDEFENDERFORIDENTITY_SERVICESOURCE
    MICROSOFTDEFENDERFORCLOUDAPPS_SERVICESOURCE
    MICROSOFTDEFENDERFOROFFICE365_SERVICESOURCE
    MICROSOFT365DEFENDER_SERVICESOURCE
    AZUREADIDENTITYPROTECTION_SERVICESOURCE
    MICROSOFTAPPGOVERNANCE_SERVICESOURCE
    DATALOSSPREVENTION_SERVICESOURCE
    UNKNOWNFUTUREVALUE_SERVICESOURCE
    MICROSOFTDEFENDERFORCLOUD_SERVICESOURCE
    MICROSOFTSENTINEL_SERVICESOURCE
    MICROSOFTINSIDERRISKMANAGEMENT_SERVICESOURCE
)

func (i ServiceSource) String() string {
    return []string{"unknown", "microsoftDefenderForEndpoint", "microsoftDefenderForIdentity", "microsoftDefenderForCloudApps", "microsoftDefenderForOffice365", "microsoft365Defender", "azureAdIdentityProtection", "microsoftAppGovernance", "dataLossPrevention", "unknownFutureValue", "microsoftDefenderForCloud", "microsoftSentinel", "microsoftInsiderRiskManagement"}[i]
}
func ParseServiceSource(v string) (any, error) {
    result := UNKNOWN_SERVICESOURCE
    switch v {
        case "unknown":
            result = UNKNOWN_SERVICESOURCE
        case "microsoftDefenderForEndpoint":
            result = MICROSOFTDEFENDERFORENDPOINT_SERVICESOURCE
        case "microsoftDefenderForIdentity":
            result = MICROSOFTDEFENDERFORIDENTITY_SERVICESOURCE
        case "microsoftDefenderForCloudApps":
            result = MICROSOFTDEFENDERFORCLOUDAPPS_SERVICESOURCE
        case "microsoftDefenderForOffice365":
            result = MICROSOFTDEFENDERFOROFFICE365_SERVICESOURCE
        case "microsoft365Defender":
            result = MICROSOFT365DEFENDER_SERVICESOURCE
        case "azureAdIdentityProtection":
            result = AZUREADIDENTITYPROTECTION_SERVICESOURCE
        case "microsoftAppGovernance":
            result = MICROSOFTAPPGOVERNANCE_SERVICESOURCE
        case "dataLossPrevention":
            result = DATALOSSPREVENTION_SERVICESOURCE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICESOURCE
        case "microsoftDefenderForCloud":
            result = MICROSOFTDEFENDERFORCLOUD_SERVICESOURCE
        case "microsoftSentinel":
            result = MICROSOFTSENTINEL_SERVICESOURCE
        case "microsoftInsiderRiskManagement":
            result = MICROSOFTINSIDERRISKMANAGEMENT_SERVICESOURCE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServiceSource(values []ServiceSource) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServiceSource) isMultiValue() bool {
    return false
}
