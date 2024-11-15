package models
type CloudPcDeviceImageStatus int

const (
    PENDING_CLOUDPCDEVICEIMAGESTATUS CloudPcDeviceImageStatus = iota
    READY_CLOUDPCDEVICEIMAGESTATUS
    FAILED_CLOUDPCDEVICEIMAGESTATUS
    UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGESTATUS
)

func (i CloudPcDeviceImageStatus) String() string {
    return []string{"pending", "ready", "failed", "unknownFutureValue"}[i]
}
func ParseCloudPcDeviceImageStatus(v string) (any, error) {
    result := PENDING_CLOUDPCDEVICEIMAGESTATUS
    switch v {
        case "pending":
            result = PENDING_CLOUDPCDEVICEIMAGESTATUS
        case "ready":
            result = READY_CLOUDPCDEVICEIMAGESTATUS
        case "failed":
            result = FAILED_CLOUDPCDEVICEIMAGESTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CLOUDPCDEVICEIMAGESTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeCloudPcDeviceImageStatus(values []CloudPcDeviceImageStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i CloudPcDeviceImageStatus) isMultiValue() bool {
    return false
}
