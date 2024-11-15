package callrecords
type ProductFamily int

const (
    UNKNOWN_PRODUCTFAMILY ProductFamily = iota
    TEAMS_PRODUCTFAMILY
    SKYPEFORBUSINESS_PRODUCTFAMILY
    LYNC_PRODUCTFAMILY
    UNKNOWNFUTUREVALUE_PRODUCTFAMILY
    AZURECOMMUNICATIONSERVICES_PRODUCTFAMILY
)

func (i ProductFamily) String() string {
    return []string{"unknown", "teams", "skypeForBusiness", "lync", "unknownFutureValue", "azureCommunicationServices"}[i]
}
func ParseProductFamily(v string) (any, error) {
    result := UNKNOWN_PRODUCTFAMILY
    switch v {
        case "unknown":
            result = UNKNOWN_PRODUCTFAMILY
        case "teams":
            result = TEAMS_PRODUCTFAMILY
        case "skypeForBusiness":
            result = SKYPEFORBUSINESS_PRODUCTFAMILY
        case "lync":
            result = LYNC_PRODUCTFAMILY
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRODUCTFAMILY
        case "azureCommunicationServices":
            result = AZURECOMMUNICATIONSERVICES_PRODUCTFAMILY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeProductFamily(values []ProductFamily) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ProductFamily) isMultiValue() bool {
    return false
}
