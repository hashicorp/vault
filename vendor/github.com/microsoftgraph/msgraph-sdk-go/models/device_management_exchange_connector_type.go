package models
// The type of Exchange Connector.
type DeviceManagementExchangeConnectorType int

const (
    // Connects to on-premises Exchange Environment.
    ONPREMISES_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE DeviceManagementExchangeConnectorType = iota
    // Connects to O365 multi-tenant Exchange environment
    HOSTED_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
    // Intune Service connects directly to O365 multi-tenant Exchange environment
    SERVICETOSERVICE_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
    // Connects to O365 Dedicated Exchange environment.
    DEDICATED_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
)

func (i DeviceManagementExchangeConnectorType) String() string {
    return []string{"onPremises", "hosted", "serviceToService", "dedicated", "unknownFutureValue"}[i]
}
func ParseDeviceManagementExchangeConnectorType(v string) (any, error) {
    result := ONPREMISES_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
    switch v {
        case "onPremises":
            result = ONPREMISES_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
        case "hosted":
            result = HOSTED_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
        case "serviceToService":
            result = SERVICETOSERVICE_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
        case "dedicated":
            result = DEDICATED_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTEXCHANGECONNECTORTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementExchangeConnectorType(values []DeviceManagementExchangeConnectorType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementExchangeConnectorType) isMultiValue() bool {
    return false
}
