package models
type CloudPcProvisioningPolicyImageType int

const (
    GALLERY_CLOUDPCPROVISIONINGPOLICYIMAGETYPE CloudPcProvisioningPolicyImageType = iota
    CUSTOM_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
    UNKNOWNFUTUREVALUE_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
)

func (i CloudPcProvisioningPolicyImageType) String() string {
    return []string{"gallery", "custom", "unknownFutureValue"}[i]
}
func ParseCloudPcProvisioningPolicyImageType(v string) (any, error) {
    result := GALLERY_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
    switch v {
        case "gallery":
            result = GALLERY_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
        case "custom":
            result = CUSTOM_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCPROVISIONINGPOLICYIMAGETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcProvisioningPolicyImageType(values []CloudPcProvisioningPolicyImageType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcProvisioningPolicyImageType) isMultiValue() bool {
    return false
}
