package models
// Device Exchange Access State.
type DeviceManagementExchangeAccessState int

const (
    // No access state discovered from Exchange
    NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATE DeviceManagementExchangeAccessState = iota
    // Device access state to Exchange is unknown
    UNKNOWN_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
    // Device has access to Exchange
    ALLOWED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
    // Device is Blocked in Exchange
    BLOCKED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
    // Device is Quarantined in Exchange
    QUARANTINED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
)

func (i DeviceManagementExchangeAccessState) String() string {
    return []string{"none", "unknown", "allowed", "blocked", "quarantined"}[i]
}
func ParseDeviceManagementExchangeAccessState(v string) (any, error) {
    result := NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
    switch v {
        case "none":
            result = NONE_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
        case "unknown":
            result = UNKNOWN_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
        case "allowed":
            result = ALLOWED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
        case "blocked":
            result = BLOCKED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
        case "quarantined":
            result = QUARANTINED_DEVICEMANAGEMENTEXCHANGEACCESSSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementExchangeAccessState(values []DeviceManagementExchangeAccessState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementExchangeAccessState) isMultiValue() bool {
    return false
}
