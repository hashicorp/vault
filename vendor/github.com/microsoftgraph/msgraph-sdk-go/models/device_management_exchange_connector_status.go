package models
// The current status of the Exchange Connector.
type DeviceManagementExchangeConnectorStatus int

const (
    // No Connector exists.
    NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS DeviceManagementExchangeConnectorStatus = iota
    // Pending Connection to the Exchange Environment.
    CONNECTIONPENDING_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    // Connected to the Exchange Environment
    CONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    // Disconnected from the Exchange Environment
    DISCONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
)

func (i DeviceManagementExchangeConnectorStatus) String() string {
    return []string{"none", "connectionPending", "connected", "disconnected", "unknownFutureValue"}[i]
}
func ParseDeviceManagementExchangeConnectorStatus(v string) (any, error) {
    result := NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
    switch v {
        case "none":
            result = NONE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "connectionPending":
            result = CONNECTIONPENDING_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "connected":
            result = CONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "disconnected":
            result = DISCONNECTED_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTEXCHANGECONNECTORSTATUS
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementExchangeConnectorStatus(values []DeviceManagementExchangeConnectorStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementExchangeConnectorStatus) isMultiValue() bool {
    return false
}
