package models
type CloudPcDomainJoinType int

const (
    AZUREADJOIN_CLOUDPCDOMAINJOINTYPE CloudPcDomainJoinType = iota
    HYBRIDAZUREADJOIN_CLOUDPCDOMAINJOINTYPE
    UNKNOWNFUTUREVALUE_CLOUDPCDOMAINJOINTYPE
)

func (i CloudPcDomainJoinType) String() string {
    return []string{"azureADJoin", "hybridAzureADJoin", "unknownFutureValue"}[i]
}
func ParseCloudPcDomainJoinType(v string) (any, error) {
    result := AZUREADJOIN_CLOUDPCDOMAINJOINTYPE
    switch v {
        case "azureADJoin":
            result = AZUREADJOIN_CLOUDPCDOMAINJOINTYPE
        case "hybridAzureADJoin":
            result = HYBRIDAZUREADJOIN_CLOUDPCDOMAINJOINTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCDOMAINJOINTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcDomainJoinType(values []CloudPcDomainJoinType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcDomainJoinType) isMultiValue() bool {
    return false
}
