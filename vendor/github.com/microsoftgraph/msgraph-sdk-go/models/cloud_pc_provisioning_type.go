package models
type CloudPcProvisioningType int

const (
    DEDICATED_CLOUDPCPROVISIONINGTYPE CloudPcProvisioningType = iota
    SHARED_CLOUDPCPROVISIONINGTYPE
    UNKNOWNFUTUREVALUE_CLOUDPCPROVISIONINGTYPE
)

func (i CloudPcProvisioningType) String() string {
    return []string{"dedicated", "shared", "unknownFutureValue"}[i]
}
func ParseCloudPcProvisioningType(v string) (any, error) {
    result := DEDICATED_CLOUDPCPROVISIONINGTYPE
    switch v {
        case "dedicated":
            result = DEDICATED_CLOUDPCPROVISIONINGTYPE
        case "shared":
            result = SHARED_CLOUDPCPROVISIONINGTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCPROVISIONINGTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcProvisioningType(values []CloudPcProvisioningType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcProvisioningType) isMultiValue() bool {
    return false
}
