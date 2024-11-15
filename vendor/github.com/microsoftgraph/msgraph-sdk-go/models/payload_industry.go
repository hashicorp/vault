package models
type PayloadIndustry int

const (
    UNKNOWN_PAYLOADINDUSTRY PayloadIndustry = iota
    OTHER_PAYLOADINDUSTRY
    BANKING_PAYLOADINDUSTRY
    BUSINESSSERVICES_PAYLOADINDUSTRY
    CONSUMERSERVICES_PAYLOADINDUSTRY
    EDUCATION_PAYLOADINDUSTRY
    ENERGY_PAYLOADINDUSTRY
    CONSTRUCTION_PAYLOADINDUSTRY
    CONSULTING_PAYLOADINDUSTRY
    FINANCIALSERVICES_PAYLOADINDUSTRY
    GOVERNMENT_PAYLOADINDUSTRY
    HOSPITALITY_PAYLOADINDUSTRY
    INSURANCE_PAYLOADINDUSTRY
    LEGAL_PAYLOADINDUSTRY
    COURIERSERVICES_PAYLOADINDUSTRY
    IT_PAYLOADINDUSTRY
    HEALTHCARE_PAYLOADINDUSTRY
    MANUFACTURING_PAYLOADINDUSTRY
    RETAIL_PAYLOADINDUSTRY
    TELECOM_PAYLOADINDUSTRY
    REALESTATE_PAYLOADINDUSTRY
    UNKNOWNFUTUREVALUE_PAYLOADINDUSTRY
)

func (i PayloadIndustry) String() string {
    return []string{"unknown", "other", "banking", "businessServices", "consumerServices", "education", "energy", "construction", "consulting", "financialServices", "government", "hospitality", "insurance", "legal", "courierServices", "IT", "healthcare", "manufacturing", "retail", "telecom", "realEstate", "unknownFutureValue"}[i]
}
func ParsePayloadIndustry(v string) (any, error) {
    result := UNKNOWN_PAYLOADINDUSTRY
    switch v {
        case "unknown":
            result = UNKNOWN_PAYLOADINDUSTRY
        case "other":
            result = OTHER_PAYLOADINDUSTRY
        case "banking":
            result = BANKING_PAYLOADINDUSTRY
        case "businessServices":
            result = BUSINESSSERVICES_PAYLOADINDUSTRY
        case "consumerServices":
            result = CONSUMERSERVICES_PAYLOADINDUSTRY
        case "education":
            result = EDUCATION_PAYLOADINDUSTRY
        case "energy":
            result = ENERGY_PAYLOADINDUSTRY
        case "construction":
            result = CONSTRUCTION_PAYLOADINDUSTRY
        case "consulting":
            result = CONSULTING_PAYLOADINDUSTRY
        case "financialServices":
            result = FINANCIALSERVICES_PAYLOADINDUSTRY
        case "government":
            result = GOVERNMENT_PAYLOADINDUSTRY
        case "hospitality":
            result = HOSPITALITY_PAYLOADINDUSTRY
        case "insurance":
            result = INSURANCE_PAYLOADINDUSTRY
        case "legal":
            result = LEGAL_PAYLOADINDUSTRY
        case "courierServices":
            result = COURIERSERVICES_PAYLOADINDUSTRY
        case "IT":
            result = IT_PAYLOADINDUSTRY
        case "healthcare":
            result = HEALTHCARE_PAYLOADINDUSTRY
        case "manufacturing":
            result = MANUFACTURING_PAYLOADINDUSTRY
        case "retail":
            result = RETAIL_PAYLOADINDUSTRY
        case "telecom":
            result = TELECOM_PAYLOADINDUSTRY
        case "realEstate":
            result = REALESTATE_PAYLOADINDUSTRY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PAYLOADINDUSTRY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePayloadIndustry(values []PayloadIndustry) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PayloadIndustry) isMultiValue() bool {
    return false
}
