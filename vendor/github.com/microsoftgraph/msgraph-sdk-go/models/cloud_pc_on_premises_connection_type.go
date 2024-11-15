package models
type CloudPcOnPremisesConnectionType int

const (
    HYBRIDAZUREADJOIN_CLOUDPCONPREMISESCONNECTIONTYPE CloudPcOnPremisesConnectionType = iota
    AZUREADJOIN_CLOUDPCONPREMISESCONNECTIONTYPE
    UNKNOWNFUTUREVALUE_CLOUDPCONPREMISESCONNECTIONTYPE
)

func (i CloudPcOnPremisesConnectionType) String() string {
    return []string{"hybridAzureADJoin", "azureADJoin", "unknownFutureValue"}[i]
}
func ParseCloudPcOnPremisesConnectionType(v string) (any, error) {
    result := HYBRIDAZUREADJOIN_CLOUDPCONPREMISESCONNECTIONTYPE
    switch v {
        case "hybridAzureADJoin":
            result = HYBRIDAZUREADJOIN_CLOUDPCONPREMISESCONNECTIONTYPE
        case "azureADJoin":
            result = AZUREADJOIN_CLOUDPCONPREMISESCONNECTIONTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCONPREMISESCONNECTIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcOnPremisesConnectionType(values []CloudPcOnPremisesConnectionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcOnPremisesConnectionType) isMultiValue() bool {
    return false
}
