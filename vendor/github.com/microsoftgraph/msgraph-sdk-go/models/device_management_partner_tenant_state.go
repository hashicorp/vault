package models
// Partner state of this tenant.
type DeviceManagementPartnerTenantState int

const (
    // Partner state is unknown.
    UNKNOWN_DEVICEMANAGEMENTPARTNERTENANTSTATE DeviceManagementPartnerTenantState = iota
    // Partner is unavailable.
    UNAVAILABLE_DEVICEMANAGEMENTPARTNERTENANTSTATE
    // Partner is enabled.
    ENABLED_DEVICEMANAGEMENTPARTNERTENANTSTATE
    // Partner connection is terminated.
    TERMINATED_DEVICEMANAGEMENTPARTNERTENANTSTATE
    // Partner messages are rejected.
    REJECTED_DEVICEMANAGEMENTPARTNERTENANTSTATE
    // Partner is unresponsive.
    UNRESPONSIVE_DEVICEMANAGEMENTPARTNERTENANTSTATE
)

func (i DeviceManagementPartnerTenantState) String() string {
    return []string{"unknown", "unavailable", "enabled", "terminated", "rejected", "unresponsive"}[i]
}
func ParseDeviceManagementPartnerTenantState(v string) (any, error) {
    result := UNKNOWN_DEVICEMANAGEMENTPARTNERTENANTSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_DEVICEMANAGEMENTPARTNERTENANTSTATE
        case "unavailable":
            result = UNAVAILABLE_DEVICEMANAGEMENTPARTNERTENANTSTATE
        case "enabled":
            result = ENABLED_DEVICEMANAGEMENTPARTNERTENANTSTATE
        case "terminated":
            result = TERMINATED_DEVICEMANAGEMENTPARTNERTENANTSTATE
        case "rejected":
            result = REJECTED_DEVICEMANAGEMENTPARTNERTENANTSTATE
        case "unresponsive":
            result = UNRESPONSIVE_DEVICEMANAGEMENTPARTNERTENANTSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementPartnerTenantState(values []DeviceManagementPartnerTenantState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementPartnerTenantState) isMultiValue() bool {
    return false
}
