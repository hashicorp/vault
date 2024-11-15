package models
type CloudPcRegionGroup int

const (
    DEFAULT_CLOUDPCREGIONGROUP CloudPcRegionGroup = iota
    AUSTRALIA_CLOUDPCREGIONGROUP
    CANADA_CLOUDPCREGIONGROUP
    USCENTRAL_CLOUDPCREGIONGROUP
    USEAST_CLOUDPCREGIONGROUP
    USWEST_CLOUDPCREGIONGROUP
    FRANCE_CLOUDPCREGIONGROUP
    GERMANY_CLOUDPCREGIONGROUP
    EUROPEUNION_CLOUDPCREGIONGROUP
    UNITEDKINGDOM_CLOUDPCREGIONGROUP
    JAPAN_CLOUDPCREGIONGROUP
    ASIA_CLOUDPCREGIONGROUP
    INDIA_CLOUDPCREGIONGROUP
    SOUTHAMERICA_CLOUDPCREGIONGROUP
    EUAP_CLOUDPCREGIONGROUP
    USGOVERNMENT_CLOUDPCREGIONGROUP
    USGOVERNMENTDOD_CLOUDPCREGIONGROUP
    NORWAY_CLOUDPCREGIONGROUP
    SWITZERLAND_CLOUDPCREGIONGROUP
    SOUTHKOREA_CLOUDPCREGIONGROUP
    UNKNOWNFUTUREVALUE_CLOUDPCREGIONGROUP
)

func (i CloudPcRegionGroup) String() string {
    return []string{"default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "norway", "switzerland", "southKorea", "unknownFutureValue"}[i]
}
func ParseCloudPcRegionGroup(v string) (any, error) {
    result := DEFAULT_CLOUDPCREGIONGROUP
    switch v {
        case "default":
            result = DEFAULT_CLOUDPCREGIONGROUP
        case "australia":
            result = AUSTRALIA_CLOUDPCREGIONGROUP
        case "canada":
            result = CANADA_CLOUDPCREGIONGROUP
        case "usCentral":
            result = USCENTRAL_CLOUDPCREGIONGROUP
        case "usEast":
            result = USEAST_CLOUDPCREGIONGROUP
        case "usWest":
            result = USWEST_CLOUDPCREGIONGROUP
        case "france":
            result = FRANCE_CLOUDPCREGIONGROUP
        case "germany":
            result = GERMANY_CLOUDPCREGIONGROUP
        case "europeUnion":
            result = EUROPEUNION_CLOUDPCREGIONGROUP
        case "unitedKingdom":
            result = UNITEDKINGDOM_CLOUDPCREGIONGROUP
        case "japan":
            result = JAPAN_CLOUDPCREGIONGROUP
        case "asia":
            result = ASIA_CLOUDPCREGIONGROUP
        case "india":
            result = INDIA_CLOUDPCREGIONGROUP
        case "southAmerica":
            result = SOUTHAMERICA_CLOUDPCREGIONGROUP
        case "euap":
            result = EUAP_CLOUDPCREGIONGROUP
        case "usGovernment":
            result = USGOVERNMENT_CLOUDPCREGIONGROUP
        case "usGovernmentDOD":
            result = USGOVERNMENTDOD_CLOUDPCREGIONGROUP
        case "norway":
            result = NORWAY_CLOUDPCREGIONGROUP
        case "switzerland":
            result = SWITZERLAND_CLOUDPCREGIONGROUP
        case "southKorea":
            result = SOUTHKOREA_CLOUDPCREGIONGROUP
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCREGIONGROUP
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcRegionGroup(values []CloudPcRegionGroup) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcRegionGroup) isMultiValue() bool {
    return false
}
