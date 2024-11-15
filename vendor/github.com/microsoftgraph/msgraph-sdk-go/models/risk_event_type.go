package models
type RiskEventType int

const (
    UNLIKELYTRAVEL_RISKEVENTTYPE RiskEventType = iota
    ANONYMIZEDIPADDRESS_RISKEVENTTYPE
    MALICIOUSIPADDRESS_RISKEVENTTYPE
    UNFAMILIARFEATURES_RISKEVENTTYPE
    MALWAREINFECTEDIPADDRESS_RISKEVENTTYPE
    SUSPICIOUSIPADDRESS_RISKEVENTTYPE
    LEAKEDCREDENTIALS_RISKEVENTTYPE
    INVESTIGATIONSTHREATINTELLIGENCE_RISKEVENTTYPE
    GENERIC_RISKEVENTTYPE
    ADMINCONFIRMEDUSERCOMPROMISED_RISKEVENTTYPE
    MCASIMPOSSIBLETRAVEL_RISKEVENTTYPE
    MCASSUSPICIOUSINBOXMANIPULATIONRULES_RISKEVENTTYPE
    INVESTIGATIONSTHREATINTELLIGENCESIGNINLINKED_RISKEVENTTYPE
    MALICIOUSIPADDRESSVALIDCREDENTIALSBLOCKEDIP_RISKEVENTTYPE
    UNKNOWNFUTUREVALUE_RISKEVENTTYPE
)

func (i RiskEventType) String() string {
    return []string{"unlikelyTravel", "anonymizedIPAddress", "maliciousIPAddress", "unfamiliarFeatures", "malwareInfectedIPAddress", "suspiciousIPAddress", "leakedCredentials", "investigationsThreatIntelligence", "generic", "adminConfirmedUserCompromised", "mcasImpossibleTravel", "mcasSuspiciousInboxManipulationRules", "investigationsThreatIntelligenceSigninLinked", "maliciousIPAddressValidCredentialsBlockedIP", "unknownFutureValue"}[i]
}
func ParseRiskEventType(v string) (any, error) {
    result := UNLIKELYTRAVEL_RISKEVENTTYPE
    switch v {
        case "unlikelyTravel":
            result = UNLIKELYTRAVEL_RISKEVENTTYPE
        case "anonymizedIPAddress":
            result = ANONYMIZEDIPADDRESS_RISKEVENTTYPE
        case "maliciousIPAddress":
            result = MALICIOUSIPADDRESS_RISKEVENTTYPE
        case "unfamiliarFeatures":
            result = UNFAMILIARFEATURES_RISKEVENTTYPE
        case "malwareInfectedIPAddress":
            result = MALWAREINFECTEDIPADDRESS_RISKEVENTTYPE
        case "suspiciousIPAddress":
            result = SUSPICIOUSIPADDRESS_RISKEVENTTYPE
        case "leakedCredentials":
            result = LEAKEDCREDENTIALS_RISKEVENTTYPE
        case "investigationsThreatIntelligence":
            result = INVESTIGATIONSTHREATINTELLIGENCE_RISKEVENTTYPE
        case "generic":
            result = GENERIC_RISKEVENTTYPE
        case "adminConfirmedUserCompromised":
            result = ADMINCONFIRMEDUSERCOMPROMISED_RISKEVENTTYPE
        case "mcasImpossibleTravel":
            result = MCASIMPOSSIBLETRAVEL_RISKEVENTTYPE
        case "mcasSuspiciousInboxManipulationRules":
            result = MCASSUSPICIOUSINBOXMANIPULATIONRULES_RISKEVENTTYPE
        case "investigationsThreatIntelligenceSigninLinked":
            result = INVESTIGATIONSTHREATINTELLIGENCESIGNINLINKED_RISKEVENTTYPE
        case "maliciousIPAddressValidCredentialsBlockedIP":
            result = MALICIOUSIPADDRESSVALIDCREDENTIALSBLOCKEDIP_RISKEVENTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_RISKEVENTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeRiskEventType(values []RiskEventType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i RiskEventType) isMultiValue() bool {
    return false
}
