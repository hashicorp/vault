package models
type CloudPcDeviceImageOsStatus int

const (
    SUPPORTED_CLOUDPCDEVICEIMAGEOSSTATUS CloudPcDeviceImageOsStatus = iota
    SUPPORTEDWITHWARNING_CLOUDPCDEVICEIMAGEOSSTATUS
    UNKNOWN_CLOUDPCDEVICEIMAGEOSSTATUS
    UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGEOSSTATUS
)

func (i CloudPcDeviceImageOsStatus) String() string {
    return []string{"supported", "supportedWithWarning", "unknown", "unknownFutureValue"}[i]
}
func ParseCloudPcDeviceImageOsStatus(v string) (any, error) {
    result := SUPPORTED_CLOUDPCDEVICEIMAGEOSSTATUS
    switch v {
        case "supported":
            result = SUPPORTED_CLOUDPCDEVICEIMAGEOSSTATUS
        case "supportedWithWarning":
            result = SUPPORTEDWITHWARNING_CLOUDPCDEVICEIMAGEOSSTATUS
        case "unknown":
            result = UNKNOWN_CLOUDPCDEVICEIMAGEOSSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGEOSSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcDeviceImageOsStatus(values []CloudPcDeviceImageOsStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcDeviceImageOsStatus) isMultiValue() bool {
    return false
}
