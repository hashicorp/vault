package security
type IoTDeviceImportanceType int

const (
    UNKNOWN_IOTDEVICEIMPORTANCETYPE IoTDeviceImportanceType = iota
    LOW_IOTDEVICEIMPORTANCETYPE
    NORMAL_IOTDEVICEIMPORTANCETYPE
    HIGH_IOTDEVICEIMPORTANCETYPE
    UNKNOWNFUTUREVALUE_IOTDEVICEIMPORTANCETYPE
)

func (i IoTDeviceImportanceType) String() string {
    return []string{"unknown", "low", "normal", "high", "unknownFutureValue"}[i]
}
func ParseIoTDeviceImportanceType(v string) (any, error) {
    result := UNKNOWN_IOTDEVICEIMPORTANCETYPE
    switch v {
        case "unknown":
            result = UNKNOWN_IOTDEVICEIMPORTANCETYPE
        case "low":
            result = LOW_IOTDEVICEIMPORTANCETYPE
        case "normal":
            result = NORMAL_IOTDEVICEIMPORTANCETYPE
        case "high":
            result = HIGH_IOTDEVICEIMPORTANCETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_IOTDEVICEIMPORTANCETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeIoTDeviceImportanceType(values []IoTDeviceImportanceType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i IoTDeviceImportanceType) isMultiValue() bool {
    return false
}
