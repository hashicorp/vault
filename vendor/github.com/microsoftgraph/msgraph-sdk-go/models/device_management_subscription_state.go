package models
// Tenant mobile device management subscription state.
type DeviceManagementSubscriptionState int

const (
    // Pending
    PENDING_DEVICEMANAGEMENTSUBSCRIPTIONSTATE DeviceManagementSubscriptionState = iota
    // Active
    ACTIVE_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    // Warning
    WARNING_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    // Disabled
    DISABLED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    // Deleted
    DELETED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    // Blocked
    BLOCKED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    // LockedOut
    LOCKEDOUT_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
)

func (i DeviceManagementSubscriptionState) String() string {
    return []string{"pending", "active", "warning", "disabled", "deleted", "blocked", "lockedOut"}[i]
}
func ParseDeviceManagementSubscriptionState(v string) (any, error) {
    result := PENDING_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
    switch v {
        case "pending":
            result = PENDING_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "active":
            result = ACTIVE_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "warning":
            result = WARNING_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "disabled":
            result = DISABLED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "deleted":
            result = DELETED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "blocked":
            result = BLOCKED_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        case "lockedOut":
            result = LOCKEDOUT_DEVICEMANAGEMENTSUBSCRIPTIONSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeDeviceManagementSubscriptionState(values []DeviceManagementSubscriptionState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i DeviceManagementSubscriptionState) isMultiValue() bool {
    return false
}
