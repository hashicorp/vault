package security
type DeviceHealthStatus int

const (
    ACTIVE_DEVICEHEALTHSTATUS DeviceHealthStatus = iota
    INACTIVE_DEVICEHEALTHSTATUS
    IMPAIREDCOMMUNICATION_DEVICEHEALTHSTATUS
    NOSENSORDATA_DEVICEHEALTHSTATUS
    NOSENSORDATAIMPAIREDCOMMUNICATION_DEVICEHEALTHSTATUS
    UNKNOWN_DEVICEHEALTHSTATUS
    UNKNOWNFUTUREVALUE_DEVICEHEALTHSTATUS
)

func (i DeviceHealthStatus) String() string {
    return []string{"active", "inactive", "impairedCommunication", "noSensorData", "noSensorDataImpairedCommunication", "unknown", "unknownFutureValue"}[i]
}
func ParseDeviceHealthStatus(v string) (any, error) {
    result := ACTIVE_DEVICEHEALTHSTATUS
    switch v {
        case "active":
            result = ACTIVE_DEVICEHEALTHSTATUS
        case "inactive":
            result = INACTIVE_DEVICEHEALTHSTATUS
        case "impairedCommunication":
            result = IMPAIREDCOMMUNICATION_DEVICEHEALTHSTATUS
        case "noSensorData":
            result = NOSENSORDATA_DEVICEHEALTHSTATUS
        case "noSensorDataImpairedCommunication":
            result = NOSENSORDATAIMPAIREDCOMMUNICATION_DEVICEHEALTHSTATUS
        case "unknown":
            result = UNKNOWN_DEVICEHEALTHSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEVICEHEALTHSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceHealthStatus(values []DeviceHealthStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceHealthStatus) isMultiValue() bool {
    return false
}
