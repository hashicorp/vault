package models
type PartnerTenantType int

const (
    MICROSOFTSUPPORT_PARTNERTENANTTYPE PartnerTenantType = iota
    SYNDICATEPARTNER_PARTNERTENANTTYPE
    BREADTHPARTNER_PARTNERTENANTTYPE
    BREADTHPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
    RESELLERPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
    VALUEADDEDRESELLERPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
    UNKNOWNFUTUREVALUE_PARTNERTENANTTYPE
)

func (i PartnerTenantType) String() string {
    return []string{"microsoftSupport", "syndicatePartner", "breadthPartner", "breadthPartnerDelegatedAdmin", "resellerPartnerDelegatedAdmin", "valueAddedResellerPartnerDelegatedAdmin", "unknownFutureValue"}[i]
}
func ParsePartnerTenantType(v string) (any, error) {
    result := MICROSOFTSUPPORT_PARTNERTENANTTYPE
    switch v {
        case "microsoftSupport":
            result = MICROSOFTSUPPORT_PARTNERTENANTTYPE
        case "syndicatePartner":
            result = SYNDICATEPARTNER_PARTNERTENANTTYPE
        case "breadthPartner":
            result = BREADTHPARTNER_PARTNERTENANTTYPE
        case "breadthPartnerDelegatedAdmin":
            result = BREADTHPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
        case "resellerPartnerDelegatedAdmin":
            result = RESELLERPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
        case "valueAddedResellerPartnerDelegatedAdmin":
            result = VALUEADDEDRESELLERPARTNERDELEGATEDADMIN_PARTNERTENANTTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PARTNERTENANTTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePartnerTenantType(values []PartnerTenantType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PartnerTenantType) isMultiValue() bool {
    return false
}
